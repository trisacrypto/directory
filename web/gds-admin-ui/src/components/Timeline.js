import classNames from 'classnames';
import React from 'react';

const Timeline = (props) => {
  const children = props.children || null;
  const Tag = props.tag;

  return (
    <Tag className={classNames('timeline-alt', 'pb-0', props.className)} {...props}>
      {children}
    </Tag>
  );
};

Timeline.defaultProps = {
  tag: 'div',
};

export default Timeline;
