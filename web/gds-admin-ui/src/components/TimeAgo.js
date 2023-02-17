import dayjs from 'dayjs';
import PropTypes from 'prop-types';
import React from 'react';
import invariant from 'tiny-invariant';

export default function TimeAgo({ time }) {
  const [, setTime] = React.useState();

  invariant(time, 'time should not be null');

  React.useEffect(() => {
    const timer = setInterval(() => {
      setTime(new Date().toLocaleString());
    }, 1000);

    return () => {
      clearInterval(timer);
    };
  }, []);

  return <span data-testid="time">{time ? dayjs(time).fromNow() : null}</span>;
}

TimeAgo.propTypes = {
  time: PropTypes.oneOfType([PropTypes.string, PropTypes.number]),
};
