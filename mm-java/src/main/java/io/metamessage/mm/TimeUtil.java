package io.metamessage.mm;

import java.time.Instant;
import java.time.LocalDate;
import java.time.LocalDateTime;
import java.time.LocalTime;
import java.time.ZoneId;
import java.time.ZoneOffset;
import java.time.ZonedDateTime;
import java.time.temporal.ChronoUnit;

final class TimeUtil {
    static final Instant EPOCH = Instant.EPOCH;

    private TimeUtil() {}

    static ZoneId zoneFromHourOffset(int hours) {
        if (hours < -12 || hours > 14) {
            throw new IllegalArgumentException("location offset hours must be between -12 and +14, got " + hours);
        }
        return ZoneOffset.ofHours(hours);
    }

    static int zoneOffsetHours(ZoneId zone) {
        if (zone == null) {
            return 0;
        }
        return ZonedDateTime.ofInstant(EPOCH, zone).getOffset().getTotalSeconds() / 3600;
    }

    static long daysSinceEpochUtc(LocalDate date) {
        return ChronoUnit.DAYS.between(LocalDate.of(1970, 1, 1), date);
    }

    static LocalDate dateFromDays(long days) {
        return LocalDate.of(1970, 1, 1).plusDays(days);
    }

    static int secondsOfDay(LocalTime t) {
        return t.toSecondOfDay();
    }

    static LocalTime timeFromSeconds(int sec) {
        return LocalTime.ofSecondOfDay(sec);
    }

    static LocalDateTime toUtcDateTime(Object o) {
        if (o instanceof LocalDateTime ldt) {
            return ldt.atZone(ZoneId.systemDefault()).withZoneSameInstant(ZoneOffset.UTC).toLocalDateTime();
        }
        if (o instanceof ZonedDateTime zdt) {
            return zdt.withZoneSameInstant(ZoneOffset.UTC).toLocalDateTime();
        }
        if (o instanceof Instant ins) {
            return LocalDateTime.ofInstant(ins, ZoneOffset.UTC);
        }
        throw new IllegalArgumentException("unsupported datetime type: " + o.getClass());
    }

    static long epochSeconds(Object o) {
        if (o instanceof LocalDateTime ldt) {
            return ldt.atZone(ZoneId.systemDefault()).toEpochSecond();
        }
        if (o instanceof ZonedDateTime zdt) {
            return zdt.toEpochSecond();
        }
        if (o instanceof Instant ins) {
            return ins.getEpochSecond();
        }
        throw new IllegalArgumentException("unsupported datetime type: " + o.getClass());
    }
}
