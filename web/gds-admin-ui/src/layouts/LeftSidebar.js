// @flow
import React, { useEffect, useRef } from 'react';
import { Link } from 'react-router-dom';
import SimpleBar from 'simplebar-react';
import AppMenu from './Menu';
import { getMenuItems } from '../helpers/menu';
import { getDirectoryLogo, getDirectoryName, getDirectoryURL } from '../utils';


const SideBarContent = ({ hideUserProfile }: SideBarContentProps) => {
    return (
        <>
            {!hideUserProfile && (
                <div className="leftbar-user">
                    <Link to="/">
                        <span className="leftbar-user-name">Dominic Keller</span>
                    </Link>
                </div>
            )}

            <AppMenu menuItems={getMenuItems()} />
        </>
    );
};


type LeftSidebarProps = {
    hideLogo: boolean,
    hideUserProfile: boolean,
    isLight: boolean,
    isCondensed: boolean,
};

const LeftSidebar = ({ isCondensed, isLight, hideLogo, hideUserProfile }: LeftSidebarProps): React$Element<any> => {
    const menuNodeRef: any = useRef(null);

    const handleOtherClick = (e: any) => {
        if (menuNodeRef && menuNodeRef.current && menuNodeRef.current.contains(e.target)) return;
        if (document.body) {
            document.body.classList.remove('sidebar-enable');
        }
    };

    useEffect(() => {
        document.addEventListener('mousedown', handleOtherClick, false);

        return () => {
            document.removeEventListener('mousedown', handleOtherClick, false);
        };
    }, []);

    return (
        <React.Fragment>
            <div className="leftside-menu" ref={menuNodeRef}>
                {!hideLogo && (
                    <React.Fragment>
                        <Link to="/" className="logo text-center logo-light">
                            <span className="logo-lg">
                                <img src={getDirectoryLogo()} alt="logo" height="36" />
                            </span>
                        </Link>
                    </React.Fragment>
                )}

                {!isCondensed && (
                    <SimpleBar style={{ maxHeight: '100%' }} timeout={500} scrollbarMaxSize={320}>
                        <SideBarContent
                            menuClickHandler={() => { }}
                            isLight={isLight}
                            hideUserProfile={hideUserProfile}
                        />
                    </SimpleBar>
                )}
                {isCondensed && <SideBarContent isLight={isLight} hideUserProfile={hideUserProfile} />}
                <ul className="side-nav position-absolute bottom-0">
                    <li className="side-nav-item">
                        <a className="side-nav-link-ref side-sub-nav-link side-nav-link" href={getDirectoryURL()}>
                            <i className="uil-exit"></i>
                            <span>Go to {getDirectoryName()}</span>
                        </a>
                    </li>
                </ul>
            </div>
        </React.Fragment>
    );
};

LeftSidebar.defaultProps = {
    hideLogo: false,
    hideUserProfile: false,
    isLight: false,
    isCondensed: false,
};

export default LeftSidebar;
