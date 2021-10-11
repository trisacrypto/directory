import React, { useState } from 'react';
import { Link } from 'react-router-dom';
import { useDispatch, useSelector } from 'react-redux';
import classNames from 'classnames';

// components
import ProfileDropdown from '../components/ProfileDropdown';
import SearchDropdown from '../components/SearchDropdown';
import TopbarSearch from '../components/TopbarSearch';

import profilePic from '../assets/images/avatar-1.jpg';
import logo from '../assets/images/gds-trisatest-logo.png';

//constants
import * as layoutConstants from '../constants/layout';
import LanguageDropdown from '../components/LanguageDropdown';
import { APICore } from '../helpers/api/apiCore';
import { loginUserSuccess } from '../redux/auth/actions';
import jwtDecode from 'jwt-decode';
import { AuthActionTypes } from '../redux/auth/constants';

// dummy search results
const SearchResults = [];

type TopbarProps = {
    hideLogo?: boolean,
    navCssClasses?: string,
    openLeftMenuCallBack?: () => void,
    topbarDark?: boolean,
};

const api = new APICore();

const Topbar = ({ hideLogo, navCssClasses, openLeftMenuCallBack, topbarDark }: TopbarProps): React$Element<any> => {
    const [isopen, setIsopen] = useState(false);
    const dispatch = useDispatch()
    const { user, loading } = useSelector(state => ({
        user: state.Auth.user,
        loading: state.Auth.loading
    }))

    React.useEffect(() => {
        const { access_token } = api.getLoggedInUser()
        const decodedToken = access_token && jwtDecode(access_token)
        dispatch(loginUserSuccess(AuthActionTypes.LOGIN_USER_SUCCESS, decodedToken))
    }, [dispatch])


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
                            {
                                !loading ? (
                                    <ProfileDropdown
                                        profilePic={user?.picture || profilePic}
                                        username={user?.name}
                                        userTitle={user?.email}
                                    />
                                ) : null
                            }
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
                                <span></span>
                                <span></span>
                                <span></span>
                            </div>
                        </Link>
                    )}

                    {/* toggle for detached layout */}
                    {layoutType === layoutConstants.LAYOUT_DETACHED && (
                        <Link to="#" className="button-menu-mobile disable-btn" onClick={handleLeftMenuCallBack}>
                            <div className="lines">
                                <span></span>
                                <span></span>
                                <span></span>
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
