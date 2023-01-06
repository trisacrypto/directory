import React from 'react';
import { Link } from 'react-router-dom';
import Select, { components } from 'react-select';
import { useHistory } from 'react-router-dom';
import classNames from 'classnames';
import { formateOptionsToLabelValueObject } from '@/utils';
import { useGetAutocompletes } from '@/hooks/useGetAutocomplete';

const TopbarSearch = (props) => {
    const [options, setOptions] = React.useState([]);
    const [inputValue, setInputValue] = React.useState('');
    const { data: autocomplete } = useGetAutocompletes();
    const [isLoading, setIsLoading] = React.useState(false);
    const [menuIsOpen, setMenuIsOpen] = React.useState(false);
    const history = useHistory();

    const handleFormatOptionLabel = (option) => {
        return (
            <Link to={`/vasps/${option.value}`} className={classNames('dropdown-item', 'notify-item', 'p-0')}>
                <span>{option.label}</span>
            </Link>
        );
    };

    const filteredOptions = (input = '') => {
        const f = formateOptionsToLabelValueObject(autocomplete);
        return f.filter(
            (option) =>
                option.value.toLowerCase().includes(input.toLowerCase()) ||
                option.label.toLowerCase().includes(input.toLowerCase())
        );
    };

    const loadOptions = (input) => {
        return new Promise((resolve, reject) => {
            if (input.length < 2) {
                return resolve([]);
            }

            setTimeout(() => {
                return resolve(filteredOptions(input));
            }, 500);
        });
    };
    const handleInputChange = async (input = '') => {
        setIsLoading(true);
        setInputValue(input);
        if (input.length < 2) {
            setMenuIsOpen(false);
            setOptions([]);
            setIsLoading(false);
        } else {
            setMenuIsOpen(true);
            const options = await loadOptions(input);
            setOptions(options);
            setIsLoading(false);
        }
    };

    const handleChange = (option, { action }) => {
        if (action === 'select-option') {
            history.push(`/vasps/${option?.value}`);
        }
    };

    return (
        <>
            <Select
                {...props}
                onBlurResetsInput
                components={{ Control, MenuList, DropdownIndicator: () => null }}
                placeholder={'Search...'}
                formatOptionLabel={handleFormatOptionLabel}
                options={options}
                value={''}
                inputValue={inputValue}
                onInputChange={handleInputChange}
                onChange={handleChange}
                menuIsOpen={menuIsOpen}
                getOptionLabel={(e) => e.label}
                maxMenuHeight="350px"
                isLoading={isLoading}
                isSearchable
                name="search-app"
                className="app-search dropdown"
                classNamePrefix="react-select"
            />
        </>
    );
};

const Control = ({ children, ...props }) => {
    const { handleClick } = props.selectProps;
    return (
        <components.Control {...props}>
            <span onMouseDown={handleClick} className="mdi mdi-magnify search-icon"></span>
            {children}
        </components.Control>
    );
};

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

export default TopbarSearch;
