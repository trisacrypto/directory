import dayjs from "dayjs";


export const isDate = (date: any) => {
    return dayjs(date).isValid();
};
