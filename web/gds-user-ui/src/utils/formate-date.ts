import dayjs from "dayjs";
export default function formatDate() {
  const dtToday = new Date();

  let month: any = dtToday.getMonth() + 1;
  let day: any = dtToday.getDate();
  const year = dtToday.getFullYear();
  if (month < 10) month = '0' + month.toString();
  if (day < 10) day = '0' + day.toString();

  const maxDate = year + '-' + month + '-' + day;

  return maxDate;
}

// format to short date with dayjs library (https://day.js.org/)
export const formatIsoDate = (date: any) => {
  // check if date arg is valid date
  if (date && dayjs(date).isValid()) {
    return dayjs(date).format('MMM D, YYYY');
  }
  return '-';
};
