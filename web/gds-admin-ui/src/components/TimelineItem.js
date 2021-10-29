import React from 'react';
import classNames from 'classnames';


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
