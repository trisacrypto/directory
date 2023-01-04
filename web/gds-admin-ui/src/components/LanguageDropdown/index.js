// @flow
import React, { useState } from 'react';
import { Dropdown } from 'react-bootstrap';
import { Link } from 'react-router-dom';

import germanyFlag from './flags/germany.jpg';
import italyFlag from './flags/italy.jpg';
import russiaFlag from './flags/russia.jpg';
import spainFlag from './flags/spain.jpg';
import enFlag from './flags/us.jpg';

// get the languages
const Languages = [
  {
    name: 'English',
    flag: enFlag,
  },
  {
    name: 'German',
    flag: germanyFlag,
  },
  {
    name: 'Italian',
    flag: italyFlag,
  },
  {
    name: 'Spanish',
    flag: spainFlag,
  },
  {
    name: 'Russian',
    flag: russiaFlag,
  },
];

const LanguageDropdown = (): React$Element<any> => {
  const enLang = Languages[0] || {};
  const [dropdownOpen, setDropdownOpen] = useState(false);

  /*
   * toggle language-dropdown
   */
  const toggleDropdown = ({ dropdownOpen: boolean }) => {
    setDropdownOpen(!dropdownOpen);
  };

  return (
    <Dropdown show={dropdownOpen} onToggle={toggleDropdown}>
      <Dropdown.Toggle
        variant="link"
        id="dropdown-languages"
        as={Link}
        to="#"
        // onClick={toggleDropdown}
        className="nav-link dropdown-toggle arrow-none"
      >
        <img src={enLang.flag} alt={enLang.name} className="me-0 me-sm-1" height="12" />{' '}
        <span className="align-middle d-none d-sm-inline-block">{enLang.name}</span>
        <i className="mdi mdi-chevron-down d-none d-sm-inline-block align-middle" />
      </Dropdown.Toggle>
      {/* <Dropdown.Menu className="dropdown-menu dropdown-menu-end dropdown-menu-animated topbar-dropdown-menu">
                <div onClick={toggleDropdown}>
                    {Languages.map((lang, i) => {
                        return (
                            <Link to="/" className="dropdown-item notify-item" key={i + '-lang'}>
                                <img src={lang.flag} alt={lang.name} className="me-1" height="12" />{' '}
                                <span className="align-middle">{lang.name}</span>
                            </Link>
                        );
                    })}
                </div>
            </Dropdown.Menu> */}
    </Dropdown>
  );
};

export default LanguageDropdown;
