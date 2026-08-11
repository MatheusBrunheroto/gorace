#!/usr/bin/env python3
"""
Faz um GET no mesmo endpoint com headers de API, em loop, ate voce cancelar (Ctrl+C).

- Respostas vazias ([] ou {}) sao ignoradas.
- Registros DUPLICADOS nao sao salvos: se um registro identico ja foi gravado
  (mesmo em execucoes anteriores), ele e ignorado.
- Intervalo ADAPTATIVO: comeca rapido (1s). Se vierem varias respostas so com
  duplicatas, desacelera (3s) para nao martelar o servidor. Assim que a resposta
  vier vazia OU trouxer algo novo, volta para o intervalo rapido.
- Gravacao APPEND-ONLY: o arquivo so cresce, nunca e reescrito/sobrescrito.
  Cada registro e um objeto JSON indentado (legivel), gravado com flush + fsync,
  entao cancelar no meio nunca corrompe nem perde o que ja foi encontrado.

Uso:
    python3 poll.py --url https://seu-site.com/api/rota \\
        -H "Authorization: Bearer TOKEN" -H "X-Api-Key: abc123"

    # headers como JSON:
    python3 poll.py --url https://... --headers '{"X-Api-Key": "abc123"}'

    # tambem le de variaveis de ambiente:
    POLL_URL=https://...  POLL_HEADERS='{"X-Api-Key":"abc"}'  python3 poll.py

Opcoes:
    --url / POLL_URL          URL do request (obrigatorio)
    -H / --header "K: V"      header (pode repetir varias vezes)
    --headers / POLL_HEADERS  headers como objeto JSON
    --out                     arquivo de saida (padrao: responses.json)
    --interval                intervalo rapido, em segundos (padrao: 1.0)
    --slow-interval           intervalo lento p/ muitas duplicatas (padrao: 3.0)
    --dupe-threshold          respostas so-duplicatas seguidas p/ desacelerar (padrao: 3)
    --timeout                 timeout de cada request em segundos (padrao: 30)
"""

import argparse
import json
import os
import sys
import time
from datetime import datetime

import requests


def parse_headers(header_list, headers_json, env_headers):
    headers = {}

    # 1) JSON vindo de env
    if env_headers:
        try:
            headers.update(json.loads(env_headers))
        except json.JSONDecodeError as e:
            sys.exit(f"POLL_HEADERS nao e um JSON valido: {e}")

    # 2) JSON vindo de --headers
    if headers_json:
        try:
            headers.update(json.loads(headers_json))
        except json.JSONDecodeError as e:
            sys.exit(f"--headers nao e um JSON valido: {e}")

    # 3) pares "Chave: Valor" vindos de -H (tem prioridade)
    for h in header_list or []:
        if ":" not in h:
            sys.exit(f'Header invalido (use "Chave: Valor"): {h!r}')
        key, _, value = h.partition(":")
        headers[key.strip()] = value.strip()

    return headers


def is_empty(data):
    """Ignora respostas vazias: [], {}, None, "" """
    if data is None:
        return True
    if isinstance(data, (list, dict, str)) and len(data) == 0:
        return True
    return False


def record_key(rec):
    """
    Chave canonica de um registro, para detectar duplicatas.
    Usa JSON com chaves ordenadas, entao a ordem dos campos nao importa.
    O campo HORARIO (adicionado por este script na gravacao) e ignorado,
    para que o mesmo registro coletado em horarios diferentes conte como duplicata.
    """
    if isinstance(rec, dict) and "HORARIO" in rec:
        rec = {k: v for k, v in rec.items() if k != "HORARIO"}
    return json.dumps(rec, sort_keys=True, ensure_ascii=False)


def load_existing(out_path):
    """
    Le o que ja foi gravado para manter a deduplicacao entre execucoes.
    O arquivo e um stream de objetos JSON (formato append-only deste script);
    tambem entende um array JSON antigo, caso o arquivo tenha sido gerado assim.
    Um objeto final incompleto (ex.: cancelado no meio de uma escrita) e ignorado
    sem quebrar a leitura. Retorna o set de chaves ja vistas.
    """
    seen = set()
    if not os.path.exists(out_path) or os.path.getsize(out_path) == 0:
        return seen

    try:
        text = open(out_path, "r", encoding="utf-8").read()
    except OSError:
        return seen

    decoder = json.JSONDecoder()
    idx, n = 0, len(text)
    while idx < n:
        while idx < n and text[idx] in " \t\r\n":
            idx += 1
        if idx >= n:
            break
        try:
            obj, end = decoder.raw_decode(text, idx)
        except ValueError:
            break  # resto incompleto/invalido: para por aqui
        idx = end
        # array antigo -> cada item e um registro; senao o proprio objeto
        for rec in (obj if isinstance(obj, list) else [obj]):
            seen.add(record_key(rec))
    return seen


def append_record(out_path, rec):
    """
    Acrescenta UM registro ao final do arquivo (append puro, nunca reescreve).
    Objeto JSON indentado + linha em branco de separador. flush + fsync para
    ser seguro contra cancelamento. Inclui o campo HORARIO com o momento
    (local, com timezone) em que o registro foi coletado.
    """
    if isinstance(rec, dict):
        rec = {**rec, "HORARIO": datetime.now().astimezone().isoformat(timespec="seconds")}
    with open(out_path, "a", encoding="utf-8") as f:
        f.write(json.dumps(rec, ensure_ascii=False, indent=2))
        f.write("\n\n")
        f.flush()
        os.fsync(f.fileno())


def main():
    parser = argparse.ArgumentParser(description="Poll adaptativo de um endpoint GET, sem duplicatas.")
    parser.add_argument("--url", default=os.environ.get("POLL_URL"),
                        help="URL do request (ou env POLL_URL)")
    parser.add_argument("-H", "--header", action="append", metavar="K: V",
                        help='Header no formato "Chave: Valor" (pode repetir)')
    parser.add_argument("--headers", default=None,
                        help="Headers como objeto JSON (ou env POLL_HEADERS)")
    parser.add_argument("--out", default="responses.json",
                        help="Arquivo de saida (padrao: responses.json)")
    parser.add_argument("--interval", type=float, default=1.0,
                        help="Intervalo rapido, em segundos (padrao: 1.0)")
    parser.add_argument("--slow-interval", type=float, default=3.0,
                        help="Intervalo lento p/ muitas duplicatas (padrao: 3.0)")
    parser.add_argument("--dupe-threshold", type=int, default=3,
                        help="Respostas so-duplicatas seguidas p/ desacelerar (padrao: 3)")
    parser.add_argument("--timeout", type=float, default=30.0,
                        help="Timeout por request em segundos (padrao: 30)")
    args = parser.parse_args()

    if not args.url:
        sys.exit("Faltou a URL. Use --url ou a variavel de ambiente POLL_URL.")

    headers = parse_headers(args.header, args.headers,
                            os.environ.get("POLL_HEADERS"))

    seen = load_existing(args.out)

    print(f"URL:       {args.url}")
    print(f"Headers:   {list(headers.keys()) or '(nenhum)'}")
    print(f"Saida:     {args.out}  (append-only)")
    print(f"Intervalo: {args.interval}s rapido / {args.slow_interval}s lento apos "
          f"{args.dupe_threshold} respostas so-duplicatas   (Ctrl+C para parar)")
    if seen:
        print(f"Continuando: {len(seen)} registro(s) ja conhecidos")
    print("-" * 50)

    session = requests.Session()
    session.headers.update(headers)

    polls = 0
    total_new = 0
    consecutive_dupes = 0
    interval = args.interval
    try:
        while True:
            start = time.monotonic()
            polls += 1
            try:
                resp = session.get(args.url, timeout=args.timeout)
                if resp.status_code != 200:
                    print(f"[#{polls}] HTTP {resp.status_code} — ignorado")
                else:
                    try:
                        data = resp.json()
                    except ValueError:
                        print(f"[#{polls}] resposta nao e JSON — ignorado")
                        data = None

                    if is_empty(data):
                        # vazio: volta ao ritmo rapido
                        if interval != args.interval:
                            print(f"[#{polls}] vazio — voltando a {args.interval}s")
                        consecutive_dupes = 0
                        interval = args.interval
                    else:
                        incoming = data if isinstance(data, list) else [data]
                        new_count = 0
                        for rec in incoming:
                            key = record_key(rec)
                            if key in seen:
                                continue  # duplicata: nao salva
                            seen.add(key)
                            append_record(args.out, rec)  # append imediato
                            new_count += 1

                        if new_count:
                            total_new += new_count
                            if interval != args.interval:
                                print(f"[#{polls}] +{new_count} novo(s) (total: {total_new}) — voltando a {args.interval}s")
                            else:
                                print(f"[#{polls}] +{new_count} novo(s)  (total: {total_new})")
                            consecutive_dupes = 0
                            interval = args.interval
                        else:
                            consecutive_dupes += 1
                            if consecutive_dupes >= args.dupe_threshold and interval != args.slow_interval:
                                interval = args.slow_interval
                                print(f"[#{polls}] so duplicatas x{consecutive_dupes} — desacelerando p/ {args.slow_interval}s")
                            else:
                                print(f"[#{polls}] so duplicatas — nada novo")
            except requests.RequestException as e:
                # erro de rede/timeout: nao derruba o loop, mantem o ritmo atual
                print(f"[#{polls}] erro de request: {e}")

            # dorme o resto do intervalo (descontando o tempo gasto no request)
            elapsed = time.monotonic() - start
            sleep_for = interval - elapsed
            if sleep_for > 0:
                time.sleep(sleep_for)
    except KeyboardInterrupt:
        print(f"\nCancelado. {polls} requests, {total_new} novos registros gravados em {args.out}")


if __name__ == "__main__":
    main()
