package io.github.metamessage.mm

import java.time.Instant
import java.time.LocalDate
import java.time.LocalDateTime
import java.time.LocalTime
import java.time.ZoneId
import java.time.ZoneOffset
import java.time.ZonedDateTime
import java.time.temporal.ChronoUnit

object TimeUtil {
    val EPOCH: Instant = Instant.EPOCH

    fun zoneFromHourOffset(hours: Int): ZoneId {
        require(hours in -12..14) { "location offset hours must be between -12 and +14, got $hours" }
        return ZoneOffset.ofHours(hours)
    }

    fun zoneOffsetHours(zone: ZoneId?): Int {
        if (zone == null) return 0
        return ZonedDateTime.ofInstant(EPOCH, zone).offset.totalSeconds / 3600
    }

    fun daysSinceEpochUtc(date: LocalDate): Long {
        return ChronoUnit.DAYS.between(LocalDate.of(1970, 1, 1), date)
    }

    fun dateFromDays(days: Long): LocalDate {
        return LocalDate.of(1970, 1, 1).plusDays(days)
    }

    fun secondsOfDay(t: LocalTime): Int {
        return t.toSecondOfDay()
    }

    fun timeFromSeconds(sec: Int): LocalTime {
        return LocalTime.ofSecondOfDay(sec.toLong())
    }

    fun toUtcDateTime(o: Any): LocalDateTime {
        return when (o) {
            is LocalDateTime -> o.atZone(ZoneId.systemDefault()).withZoneSameInstant(ZoneOffset.UTC).toLocalDateTime()
            is ZonedDateTime -> o.withZoneSameInstant(ZoneOffset.UTC).toLocalDateTime()
            is Instant -> LocalDateTime.ofInstant(o, ZoneOffset.UTC)
            else -> throw IllegalArgumentException("unsupported datetime type: ${o.javaClass}")
        }
    }

    fun epochSeconds(o: Any): Long {
        return when (o) {
            is LocalDateTime -> o.atZone(ZoneId.systemDefault()).toEpochSecond()
            is ZonedDateTime -> o.toEpochSecond()
            is Instant -> o.epochSecond
            else -> throw IllegalArgumentException("unsupported datetime type: ${o.javaClass}")
        }
    }
}
