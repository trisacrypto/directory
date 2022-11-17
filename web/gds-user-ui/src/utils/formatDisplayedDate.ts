import dayjs from 'dayjs';

const DATE_FORMAT = 'DD-MM-YYYY';
const formatDisplayedDate = (
  date: string | number | Date | dayjs.Dayjs,
  format = DATE_FORMAT
): string | 'N/A' => {
  if (dayjs(date).isValid()) {
    return dayjs(date).format(format);
  }

  return 'N/A';
};

export default formatDisplayedDate;
