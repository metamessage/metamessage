import { MMEncoder } from './src/mm/encoder.js';
import { intLen, MMPrefix, MMConstants } from './src/mm/constants.js';

function traceEncode(v, isNegative) {
    const encoder = new MMEncoder();
    encoder.encodeInt(v);
    const data = encoder.buffer.data;
    const bytes = Array.from(data).map(b => '0x' + b.toString(16).padStart(2, '0'));

    console.log('\n=== Encoding ' + v + ' ===');
    console.log('Bytes:', bytes.join(' '));
    console.log('Total length:', data.length);

    // Trace through intLen
    const firstByte = data[0];
    const lenResult = intLen(firstByte);
    console.log('First byte:', '0x' + firstByte.toString(16), '=', firstByte);
    console.log('intLen result:', lenResult);

    // Simulate decode
    let offset = 1;
    let value = 0n;
    for (let i = 0; i < lenResult.extraBytes; i++) {
        const b = data[offset++];
        value = (value << 8n) | BigInt(b);
    }
    const suffix = firstByte & 0x1F;
    value = (value << 8n) | BigInt(suffix);

    if (isNegative) {
        value = -value;
    }

    console.log('Decoded value:', value);
    console.log('Expected:', v);
    console.log('Match:', value === v ? 'YES' : 'NO - MISMATCH!');
}

traceEncode(0n, false);
traceEncode(23n, false);
traceEncode(255n, false);
traceEncode(256n, false);
traceEncode(65535n, false);
traceEncode(65536n, false);
traceEncode(123456n, false);
traceEncode(-1n, true);
traceEncode(-7890n, true);
traceEncode(-9223372036854775808n, true); // MIN_INT64