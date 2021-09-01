// @flow
import React from 'react';
import { Link } from 'react-router-dom';
import Select, { components } from 'react-select';
import classNames from 'classnames';

/*
 * get options
 */
const optionGetter = (option) => {
    switch (option.type) {
        case 'report':
            return (
                <Link to="/" className={classNames('dropdown-item', 'notify-item', 'p-0')}>
                    <i className={classNames(option.icon, 'font-16', 'me-1')}></i>
                    <span>{option.label}</span>
                </Link>
            );
        case 'help':
            return (
                <Link to="/" className={classNames('dropdown-item', 'notify-item', 'p-0')}>
                    <i className={classNames(option.icon, 'font-16', 'me-1')}></i>
                    <span>{option.label}</span>
                </Link>
            );
        case 'settings':
            return (
                <Link to="/" className={classNames('dropdown-item', 'notify-item', 'p-0')}>
                    <i className={classNames(option.icon, 'font-16', 'me-1')}></i>
                    <span>{option.label}</span>
                </Link>
            );
        case 'title':
            return (
                <div className="noti-title">
                    <h6 className="text-overflow mb-2 text-uppercase">Users</h6>
                </div>
            );
        case 'users':
            return (
                <>
                    <Link to="/" className="dropdown-item notify-item p-0">
                        <div className="d-flex">
                            <img
                                src={option.userDetails.avatar}
                                alt=""
                                className="d-flex me-2 rounded-circle"
                                height="32"
                            />
                            <div className="w-100">
                                <h5 className="m-0 font-14">
                                    {option.userDetails.firstname} {option.userDetails.lastname}
                                </h5>
                                <span className="font-12 mb-0">{option.userDetails.position}</span>
                            </div>
                        </div>
                    </Link>
                </>
            );

        default:
            return;
    }
};


/* custon control */
const Control = ({ children, ...props }) => {
    const { handleClick } = props.selectProps;
    return (
        <components.Control {...props}>
            <span onMouseDown={handleClick} className="mdi mdi-magnify search-icon"></span>
            {children}
        </components.Control>
    );
};

/* custon indicator */
const IndicatorsContainer = (props) => {
    const { handleClick } = props.selectProps;
    return (
        <div style={{ }}>
            <components.IndicatorsContainer {...props}>
                <button className="btn btn-primary" onMouseDown={handleClick}>
                    Search
                </button>
            </components.IndicatorsContainer>
        </div>
    );
};

/* custom menu list */
const MenuList = (props) => {
    const { options } = props.selectProps;

    return (
        <components.MenuList {...props}>
            {/* menu header */}
            <div className="dropdown-header noti-title">
                <h5 className="text-overflow mb-2">
                    Found <span className="text-danger">{options.length}</span> results
                </h5>
            </div>
            {props.children}
        </components.MenuList>
    );
};

/* fomates the option label */
const handleFormatOptionLabel = (option) => {
    const formattedOption = optionGetter(option);
    return <div>{formattedOption}</div>;
};

type SearchResultItem = {
    id: number,
    title: string,
    redirectTo: string,
    icon: string,
};

type TopbarSearchProps = {
    items: Array<SearchResultItem>,
};

const TopbarSearch = (props: TopbarSearchProps): React$Element<any> => {

    const onClick = (e) => {
        e.preventDefault();
        e.stopPropagation();
    };

    return (
        <>
            <Select
                {...props}
                components={{ Control, IndicatorsContainer, MenuList }}
                placeholder={'Search...'}
                formatOptionLabel={handleFormatOptionLabel}
                isOptionDisabled={(option) => option.type === 'title'}
                maxMenuHeight="350px"
                handleClick={onClick}
                isSearchable
                name="search-app"
                className="app-search dropdown"
                classNamePrefix="react-select"
            />
        </>
    );
};

export default TopbarSearch;
