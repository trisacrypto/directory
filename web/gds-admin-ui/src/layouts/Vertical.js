import React, { Suspense, useCallback, useEffect, useState } from 'react';
import { Container } from 'react-bootstrap';
import { useDispatch, useSelector } from 'react-redux';

import * as layoutConstants from '../constants/layout';
// actions
import { changeSidebarTheme, changeSidebarType } from '../redux/actions';

const Topbar = React.lazy(() => import('./Topbar'));
const LeftSidebar = React.lazy(() => import(/* webpackPreload: true */ './LeftSidebar'));
const Footer = React.lazy(() => import(/* webpackPreload: true */ './Footer'));

const loading = () => <div className="" />;

const VerticalLayout = ({ children }, state) => {
    const dispatch = useDispatch();

    const { leftSideBarTheme, leftSideBarType } = useSelector((state) => ({
        layoutWidth: state.Layout.layoutWidth,
        leftSideBarTheme: state.Layout.leftSideBarTheme,
        leftSideBarType: state.Layout.leftSideBarType,
        showRightSidebar: state.Layout.showRightSidebar,
    }));

    const [isMenuOpened, setIsMenuOpened] = useState(false);

    /**
     * Open the menu when having mobile screen
     */
    const openMenu = () => {
        setIsMenuOpened((prevState) => {
            setIsMenuOpened(!prevState);
        });

        if (document.body) {
            if (isMenuOpened) {
                document.body.classList.remove('sidebar-enable');
            } else {
                document.body.classList.add('sidebar-enable');
            }
        }
    };

    const updateDimensions = useCallback(() => {
        // activate the condensed sidebar if smaller devices like ipad or tablet
        if (window.innerWidth >= 768 && window.innerWidth <= 1028) {
            dispatch(changeSidebarType(layoutConstants.LEFT_SIDEBAR_TYPE_CONDENSED));
        } else if (window.innerWidth > 1028) {
            dispatch(changeSidebarType(layoutConstants.LEFT_SIDEBAR_TYPE_FIXED));
        }
    }, [dispatch]);

    useEffect(() => {
        dispatch(changeSidebarTheme(layoutConstants.LEFT_SIDEBAR_THEME_DARK));

        // activate the condensed sidebar if smaller devices like ipad or tablet
        if (window.innerWidth >= 768 && window.innerWidth <= 1028) {
            dispatch(changeSidebarType(layoutConstants.LEFT_SIDEBAR_TYPE_CONDENSED));
        }

        window.addEventListener('resize', updateDimensions);

        return () => {
            window.removeEventListener('resize', updateDimensions);
        };
    }, [dispatch, updateDimensions]);

    const isCondensed = leftSideBarType === layoutConstants.LEFT_SIDEBAR_TYPE_CONDENSED;
    const isLight = leftSideBarTheme === layoutConstants.LEFT_SIDEBAR_THEME_LIGHT;

    return (
        <div className="wrapper">
            <LeftSidebar isCondensed={isCondensed} isLight={isLight} hideUserProfile={true} />
            <div className="content-page">
                <div className="content">
                    <Topbar openLeftMenuCallBack={openMenu} hideLogo={true} />
                    <Container fluid>
                        <Suspense fallback={loading()}>{children}</Suspense>
                    </Container>
                </div>
                <Footer />
            </div>
        </div>
    );
};
export default VerticalLayout;
