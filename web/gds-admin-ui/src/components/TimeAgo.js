import React from "react";
import dayjs from "dayjs";

export default function TimeAgo({ time }) {
    const [, setTime] = React.useState();

    React.useEffect(() => {
        const timer = setInterval(() => {
            setTime(new Date().toLocaleString());
        }, 1000);

        return () => {
            clearInterval(timer);
        };
    }, []);

    return dayjs(time).fromNow()
}