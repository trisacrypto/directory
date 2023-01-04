import React, { useEffect, useRef } from 'react';
import { Link } from 'react-router-dom';
import SimpleBar from 'simplebar-react';

import { getMenuItems } from '@/helpers/menu';
import { getDirectoryLogo, getDirectoryName, getDirectoryURL } from '@/utils';

import AppMenu from './Menu';

const SideBarContent = ({ hideUserProfile }) => (
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

const LeftSidebar = ({ isCondensed, isLight, hideLogo, hideUserProfile }) => {
  const menuNodeRef = useRef(null);

  const handleOtherClick = (e) => {
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
    <div className="leftside-menu" ref={menuNodeRef}>
      {!hideLogo && (
        <Link to="/" className="logo text-center logo-light">
          <span className="logo-lg">
            <img src={getDirectoryLogo()} alt="logo" height="36" />
          </span>
        </Link>
      )}

      {!isCondensed && (
        <SimpleBar style={{ maxHeight: '100%' }} timeout={500} scrollbarMaxSize={320}>
          <SideBarContent
            menuClickHandler={() => {}}
            isLight={isLight}
            hideUserProfile={hideUserProfile}
          />
        </SimpleBar>
      )}
      {isCondensed && <SideBarContent isLight={isLight} hideUserProfile={hideUserProfile} />}
      <ul className="side-nav position-absolute bottom-0">
        <li className="side-nav-item">
          <a className="side-nav-link-ref side-sub-nav-link side-nav-link" href={getDirectoryURL()}>
            <i className="uil-exit" />
            <span>Go to {getDirectoryName()}</span>
          </a>
        </li>
      </ul>
    </div>
  );
};

LeftSidebar.defaultProps = {
  hideLogo: false,
  hideUserProfile: false,
  isLight: false,
  isCondensed: false,
};

export default LeftSidebar;
