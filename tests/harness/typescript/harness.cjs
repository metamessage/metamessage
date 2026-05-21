#!/usr/bin/env node
/**
 * MetaMessage TypeScript test harness - parse JSONC file and re-print to JSONC.
 */
const fs = require('fs');
const path = require('path');

// Resolve the mm-ts module
const mmPath = path.resolve(__dirname, '..', '..', '..', 'mm-ts');
const { parseJSONC, toJSONC } = require(path.join(mmPath, 'dist', 'jsonc', 'index'));

const args = process.argv.slice(2);
if (args.length < 1) {
    console.error('usage: harness.cjs <file.jsonc>');
    process.exit(1);
}

const data = fs.readFileSync(args[0], 'utf-8');

try {
    const doc = parseJSONC(data);
    const output = toJSONC(doc);
    process.stdout.write(output);
} catch (e) {
    console.error('parse error:', e.message);
    process.exit(1);
}