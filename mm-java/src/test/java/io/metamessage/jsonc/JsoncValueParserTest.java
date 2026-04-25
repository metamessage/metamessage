package io.metamessage.jsonc;

import static org.junit.jupiter.api.Assertions.assertArrayEquals;
import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertInstanceOf;
import static org.junit.jupiter.api.Assertions.assertThrows;
import static org.junit.jupiter.api.Assertions.assertTrue;

import io.metamessage.mm.MmTag;
import io.metamessage.mm.ValueType;
import java.math.BigInteger;
import java.net.URI;
import java.time.LocalDate;
import java.time.LocalDateTime;
import java.time.LocalTime;
import java.util.Base64;
import java.util.UUID;
import org.junit.jupiter.api.Test;

class JsoncValueParserTest {

    @Test
    void stringLiteralCoversMainScalarTypes() {
        MmTag str = MmTag.empty();
        str.type = ValueType.STRING;
        assertEquals("abc", JsoncValueParser.stringLiteral("abc", str, "p").data());

        MmTag dt = MmTag.empty();
        dt.type = ValueType.DATETIME;
        Object d1 = JsoncValueParser.stringLiteral("2026-01-02 03:04:05", dt, "p").data();
        assertEquals(LocalDateTime.of(2026, 1, 2, 3, 4, 5), d1);

        MmTag date = MmTag.empty();
        date.type = ValueType.DATE;
        assertEquals(LocalDate.of(2026, 1, 2), JsoncValueParser.stringLiteral("2026-01-02", date, "p").data());

        MmTag time = MmTag.empty();
        time.type = ValueType.TIME;
        assertEquals(LocalTime.of(3, 4, 5), JsoncValueParser.stringLiteral("03:04:05", time, "p").data());

        MmTag url = MmTag.empty();
        url.type = ValueType.URL;
        assertEquals(URI.create("https://example.com/a"), JsoncValueParser.stringLiteral("https://example.com/a", url, "p").data());

        MmTag uuid = MmTag.empty();
        uuid.type = ValueType.UUID;
        UUID u = UUID.randomUUID();
        assertEquals(u, JsoncValueParser.stringLiteral(u.toString(), uuid, "p").data());

        MmTag bytes = MmTag.empty();
        bytes.type = ValueType.BYTES;
        byte[] payload = "hello".getBytes();
        Object parsedBytes = JsoncValueParser.stringLiteral(Base64.getEncoder().encodeToString(payload), bytes, "p").data();
        assertInstanceOf(byte[].class, parsedBytes);
        assertArrayEquals(payload, (byte[]) parsedBytes);
    }

    @Test
    void stringLiteralNullAndErrorCases() {
        MmTag t = MmTag.empty();
        t.type = ValueType.STRING;
        t.isNull = true;
        assertEquals("", JsoncValueParser.stringLiteral("", t, "p").data());
        assertThrows(JsoncException.class, () -> JsoncValueParser.stringLiteral("x", t, "p"));

        MmTag bad = MmTag.empty();
        bad.type = ValueType.FLOAT32;
        assertThrows(JsoncException.class, () -> JsoncValueParser.stringLiteral("x", bad, "p"));
    }

    @Test
    void numberLiteralCoversSignedUnsignedFloatBigInt() {
        MmTag i64 = MmTag.empty();
        i64.type = ValueType.INT64;
        assertEquals(1234567890123L, JsoncValueParser.numberLiteral("1234567890123", i64, "p").data());

        MmTag u = MmTag.empty();
        u.type = ValueType.UINT;
        assertEquals(42L, JsoncValueParser.numberLiteral("42", u, "p").data());
        assertThrows(JsoncException.class, () -> JsoncValueParser.numberLiteral("-1", u, "p"));

        MmTag f32 = MmTag.empty();
        f32.type = ValueType.FLOAT32;
        assertEquals(1.25f, JsoncValueParser.numberLiteral("1.25", f32, "p").data());

        MmTag exp = MmTag.empty();
        exp.type = ValueType.FLOAT64;
        Object expVal = JsoncValueParser.numberLiteral("1e3", exp, "p").data();
        assertEquals(1000.0d, (double) expVal);

        MmTag bi = MmTag.empty();
        bi.type = ValueType.BIGINT;
        assertEquals(new BigInteger("-9999999999999999999999"), JsoncValueParser.numberLiteral("-9999999999999999999999", bi, "p").data());
    }

    @Test
    void boolLiteralTypeChecks() {
        MmTag unknown = MmTag.empty();
        JcNode.JcValue v = JsoncValueParser.boolLiteral(true, unknown, "p");
        assertEquals(true, v.data());
        assertEquals(ValueType.BOOL, v.tag().type);

        MmTag nonBool = MmTag.empty();
        nonBool.type = ValueType.STRING;
        JsoncException ex = assertThrows(JsoncException.class, () -> JsoncValueParser.boolLiteral(true, nonBool, "p"));
        assertTrue(ex.getMessage().contains("type mismatch"));
    }
}
