import React from "react"
import classNames from 'classnames';

const Spinner = (props) => {
    const children = props.children || null;
    const Tag = props.tag || 'div';
    const color = props.color || 'secondary';
    const size = props.size || '';

    return (
        <Tag
            role="status"
            className={classNames(
                { 'spinner-border': props.type === 'bordered', 'spinner-grow': props.type === 'grow' },
                [`text-${color}`],
                { [`avatar-${size}|`]: size },
                props.className
            )}>
            {children}
        </Tag>
    );
};

Spinner.defaultProps = {
    tag: 'div',
    type: 'bordered',
};

export default Spinner;