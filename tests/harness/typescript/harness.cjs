#!/usr/bin/env node
/**
 * MetaMessage TypeScript test harness - parse JSONC file and re-print to JSONC.
 */
const fs = require('fs');
const path = require('path');

// Resolve the mm-ts module
const mmPath = path.resolve(__dirname, '..', '..', '..', 'mm-ts');
const { parseJSONC, toJSONC } = require(path.join(mmPath, 'dist', 'jsonc', 'index'));
const { fromJSONC, decodeToJSONC } = require(path.join(mmPath, 'dist', 'core', 'index'));

const args = process.argv.slice(2);
if (args.length < 1) {
    console.error('usage: harness.cjs [--encode|--decode] <file.jsonc>');
    process.exit(1);
}

if (args[0] === '--encode') {
    if (args.length < 2) {
        console.error('usage: harness.cjs --encode <file.jsonc>');
        process.exit(1);
    }
    const data = fs.readFileSync(args[1], 'utf-8');
    try {
        const wire = fromJSONC(data);
        process.stdout.write(Buffer.from(wire).toString('hex'));
    } catch (e) {
        console.error('encode error:', e.message);
        process.exit(1);
    }
    return;
}

if (args[0] === '--decode') {
    const hexStr = fs.readFileSync('/dev/stdin', 'utf-8').trim();
    try {
        const wire = Buffer.from(hexStr, 'hex');
        const output = decodeToJSONC(new Uint8Array(wire));
        process.stdout.write(output);
    } catch (e) {
        console.error('decode error:', e.message);
        process.exit(1);
    }
    return;
}

// Existing behavior
const data = fs.readFileSync(args[0], 'utf-8');
try {
    const doc = parseJSONC(data);
    const output = toJSONC(doc);
    process.stdout.write(output);
} catch (e) {
    console.error('parse error:', e.message);
    process.exit(1);
}