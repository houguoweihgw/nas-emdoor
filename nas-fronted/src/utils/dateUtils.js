// dateUtils.js

// 将日期格式化为 "几年几月几日 周几 时区 几点几分几秒" 格式
function formatDateTime(dateTimeStr) {
    const date = new Date(dateTimeStr);

    const year = date.getFullYear();
    const month = date.getMonth() + 1;
    const day = date.getDate();
    const dayOfWeek = getDayOfWeek(date);
    const hours = date.getHours().toString().padStart(2, '0');
    const minutes = date.getMinutes().toString().padStart(2, '0');
    const seconds = date.getSeconds().toString().padStart(2, '0');

    // const formattedDateTime = `${year}年${month}月${day}日 周${dayOfWeek} ${getTimezone()} ${hours}:${minutes}:${seconds}`;
    const formattedDateTime = `${year}年${month}月${day}日 周${dayOfWeek} ${hours}:${minutes}:${seconds}`;
    return formattedDateTime;
}

function kbToMb(kb) {
    const mb = kb / 1000000;
    return mb.toFixed(1); // 保留一位小数
}

function getLatitude(latitude) {
    if (latitude > 0)
        return latitude.toFixed(2).toString() + "N"; // 保留一位小数
    else
        return latitude.toFixed(2).toString() + "S";
}


function getLongitude(longitude) {
    if (longitude > 0)
        return longitude.toFixed(2).toString() + "E"; // 保留一位小数
    else
        return longitude.toFixed(2).toString() + "W";
}

// 获取日期是星期几
function getDayOfWeek(date) {
    const days = ["日", "一", "二", "三", "四", "五", "六"];
    return days[date.getDay()];
}

function getAltitude(altitude) {
    return altitude.toFixed(2).toString() + "m"; // 保留一位小数
}

// 获取时区
function getTimezone() {
    const offset = new Date().getTimezoneOffset();
    const hours = Math.floor(Math.abs(offset) / 60);
    const sign = offset < 0 ? "+" : "-";

    return `UTC${sign}${hours.toString().padStart(2, '0')}`;
}

export {formatDateTime, kbToMb, getLatitude, getLongitude,getAltitude};
