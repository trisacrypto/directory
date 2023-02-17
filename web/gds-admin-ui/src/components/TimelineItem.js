import classNames from 'classnames';
import React from 'react';

const TimelineItem = (props) => {
  const children = props.children || null;
  const Tag = props.tag;

  return (
    <Tag className={classNames('timeline-item', props.className)} {...props}>
      {children}
    </Tag>
  );
};

TimelineItem.defaultProps = {
  tag: 'div',
};

export default TimelineItem;
