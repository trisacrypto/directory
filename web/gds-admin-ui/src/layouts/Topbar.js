import classNames from 'classnames';
import React, { useState } from 'react';
import { useSelector } from 'react-redux';
import { Link } from 'react-router-dom';

import profilePic from '@/assets/images/avatar-1.jpg';
import logo from '@/assets/images/gds-trisatest-logo.png';
import LanguageDropdown from '@/components/LanguageDropdown';
// components
import ProfileDropdown from '@/components/ProfileDropdown';
import SearchDropdown from '@/components/SearchDropdown';
import TopbarSearch from '@/components/TopbarSearch';
// constants
import * as layoutConstants from '@/constants/layout';
import useAuth from '@/contexts/auth/use-auth';

// dummy search results
const SearchResults = [];

const Topbar = ({ hideLogo, navCssClasses, openLeftMenuCallBack, topbarDark }) => {
    const [isopen, setIsopen] = useState(false);
    const { user } = useAuth();

    const navbarCssClasses = navCssClasses || '';
    const containerCssClasses = !hideLogo ? 'container-fluid' : '';

    const { layoutType } = useSelector((state) => ({
        layoutType: state.Layout.layoutType,
    }));

    const handleLeftMenuCallBack = () => {
        setIsopen((prevState) => !prevState);
        if (openLeftMenuCallBack) openLeftMenuCallBack();
    };

    return (
        <React.Fragment>
            <div className={`navbar-custom ${navbarCssClasses}`}>
                <div className={containerCssClasses}>
                    {!hideLogo && (
                        <Link to="/" className="topnav-logo">
                            <span className="topnav-logo-lg">
                                <img src={logo} alt="logo" height="16" />
                            </span>
                        </Link>
                    )}

                    <ul className="list-unstyled topbar-menu float-end mb-0">
                        <li className="notification-list topbar-dropdown d-xl-none">
                            <SearchDropdown />
                        </li>
                        <li className="dropdown notification-list topbar-dropdown d-none d-lg-block">
                            <LanguageDropdown />
                        </li>
                        <li className="dropdown notification-list">
                            {user ? (
                                <ProfileDropdown
                                    profilePic={user?.picture || profilePic}
                                    username={user?.name}
                                    userTitle={user?.email}
                                />
                            ) : null}
                        </li>
                    </ul>

                    {/* toggle for vertical layout */}
                    {layoutType === layoutConstants.LAYOUT_VERTICAL && (
                        <button className="button-menu-mobile open-left disable-btn" onClick={handleLeftMenuCallBack}>
                            <i className="mdi mdi-menu" />
                        </button>
                    )}

                    {/* toggle for horizontal layout */}
                    {layoutType === layoutConstants.LAYOUT_HORIZONTAL && (
                        <Link
                            to="#"
                            className={classNames('navbar-toggle', { open: isopen })}
                            onClick={handleLeftMenuCallBack}>
                            <div className="lines">
                                <span />
                                <span />
                                <span />
                            </div>
                        </Link>
                    )}

                    {/* toggle for detached layout */}
                    {layoutType === layoutConstants.LAYOUT_DETACHED && (
                        <Link to="#" className="button-menu-mobile disable-btn" onClick={handleLeftMenuCallBack}>
                            <div className="lines">
                                <span />
                                <span />
                                <span />
                            </div>
                        </Link>
                    )}
                    <TopbarSearch items={SearchResults} />
                </div>
            </div>
        </React.Fragment>
    );
};

export default Topbar;
