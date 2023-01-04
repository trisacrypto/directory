import dayjs from 'dayjs';
import isDate from 'lodash/isDate';

/**
 * Verify that the passed date is not older than 30days
 * @param {Date} pastDate
 * @returns boolean
 */
const isRecent = (pastDate = '') => {
  const _pastDate = new Date(pastDate);
  if (isDate(_pastDate)) {
    const now = dayjs();
    const THIRTY_DAYS_IN_MS = 30 * 24 * 60 * 60 * 1000; // 30days
    const timeDiffInMs = now.diff(_pastDate);

    return timeDiffInMs <= THIRTY_DAYS_IN_MS;
  }
  return false;
};

export default isRecent;
