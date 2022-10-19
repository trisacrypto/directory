import React, { useState } from 'react';
import { Link } from 'react-router-dom';
import { Dropdown } from 'react-bootstrap';
import useAuth from 'contexts/auth/use-auth';


const ProfileDropdown = (props, state) => {
    const profilePic = props.profilePic || null;
    const [dropdownOpen, setDropdownOpen] = useState(false);
    const { logout } = useAuth()

    const toggleDropdown = () => {
        setDropdownOpen(!dropdownOpen);
    };

    const handleLogoutClick = async () => {
        await logout()
    }

    return (
        <Dropdown show={dropdownOpen} onToggle={toggleDropdown}>
            <Dropdown.Toggle
                variant="link"
                id="dropdown-profile"
                as={Link}
                to="#"
                onClick={toggleDropdown}
                className="nav-link dropdown-toggle nav-user arrow-none me-0">
                <span className="account-user-avatar">
                    <img src={profilePic} className="rounded-circle" alt="user" />
                </span>
                <span>
                    <span className="account-user-name">{props.username}</span>
                    <span className="account-position">{props.userTitle}</span>
                </span>
            </Dropdown.Toggle>
            <Dropdown.Menu className="dropdown-menu dropdown-menu-end dropdown-menu-animated topbar-dropdown-menu profile-dropdown">
                <div onClick={toggleDropdown}>
                    <div className="dropdown-header noti-title">
                        <h6 className="text-overflow m-0">Welcome!</h6>
                    </div>
                    <a target="_top" href="mailto:support@rotational.io" className="dropdown-item notify-item">
                        <i className={`mdi mdi-help me-1`}></i>
                        <span>Support</span>
                    </a>
                    <a target="_blank" rel="noreferrer" href="https://trisa-workspace.slack.com" className="dropdown-item notify-item">
                        <i className={`mdi mdi-launch me-1`}></i>
                        <span>Slack</span>
                    </a>
                    <div onClick={handleLogoutClick} className="dropdown-item notify-item" style={{ cursor: 'pointer' }}>
                        <i className={`mdi mdi-logout me-1`}></i>
                        <span>Logout</span>
                    </div>
                </div>
            </Dropdown.Menu>
        </Dropdown>
    );
};

export default ProfileDropdown;
