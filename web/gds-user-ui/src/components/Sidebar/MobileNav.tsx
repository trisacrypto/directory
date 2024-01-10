import {
  Flex,
  FlexProps,
  HStack,
  IconButton,
  useColorModeValue,
  Text,
  Divider,
  Menu,
  MenuButton,
  Box,
  Avatar,
  MenuList,
  MenuItem,
  MenuDivider,
  Show,
  Tooltip
} from '@chakra-ui/react';
import { FiMenu } from 'react-icons/fi';
import LanguagesDropdown from 'components/LanguagesDropdown';
import { useDispatch, useSelector } from 'react-redux';
import { clearCookies, clearLocalStorage } from 'utils/cookies';
import { Link, useNavigate } from 'react-router-dom';
import DefaultAvatar from 'assets/default_avatar.svg';
import { resetStore } from 'application/store';
import { userSelector, logout } from 'modules/auth/login/user.slice';
import { Trans } from '@lingui/react';
import { t } from '@lingui/macro';
import { colors } from 'utils/theme';
import { APP_PATH } from 'utils/constants';
import { canCreateOrganization } from 'utils/permission';
import useMemberState from 'modules/dashboard/member/hooks/useMemberState';

interface MobileProps extends FlexProps {
  onOpen: () => void;
  isLoading?: boolean;
}
const MobileNav = ({ onOpen, ...rest }: MobileProps) => {
  const { setDefaultNetwork } = useMemberState();
  const dispatch = useDispatch();
  const { user } = useSelector(userSelector);
  console.log('user', user);
  const navigate = useNavigate();
  const handleLogout = (e: any) => {
    e.preventDefault();
    clearCookies();
    clearLocalStorage();
    localStorage.removeItem('persist:root');
    // reset store

    setDefaultNetwork();
    dispatch(logout());
    resetStore();
    navigate('/');
  };

  return (
    <Flex
      ml={{ base: 0, md: 64 }}
      px={{ base: 2, md: 4 }}
      alignItems="center"
      bg={useColorModeValue('white', 'gray.900')}
      borderBottomWidth="1px"
      borderBottomColor={useColorModeValue('gray.200', 'gray.700')}
      justifyContent={{ base: 'space-between', md: 'flex-end' }}
      {...rest}>
      <IconButton
        display={{ base: 'flex', md: 'none' }}
        onClick={onOpen}
        variant="outline"
        aria-label="open menu"
        icon={<FiMenu />}
      />
      <HStack
        spacing={{ base: 1, lg: 2, xl: 4 }}
        w={{ base: '100%', md: 'none' }}
        justifyContent="end">
        <HStack>
          <LanguagesDropdown />
        </HStack>
        <HStack>
          <Tooltip label={t`Current Organization name`} hasArrow>
            <Text fontWeight={'bold'} color={colors.system.blue}>
              {user?.vasp?.name || 'N/A'}
            </Text>
          </Tooltip>
        </HStack>
        <Show above="lg">
          <Divider orientation="vertical" height={8} />
        </Show>
        <Menu>
          <MenuButton data-testid="menu" transition="all 0.3s" _focus={{ boxShadow: 'none' }}>
            <HStack>
              <Show above="lg">
                <Text fontSize="sm" color="blackAlpha.700">
                  {user?.name || <Trans id="Guest">Guest</Trans>}
                </Text>
              </Show>
              <Box borderRadius="50%" borderWidth={2} padding={0.5}>
                <Avatar
                  size={'md'}
                  height="43.3px"
                  width="43.3px"
                  src={user?.pictureUrl || DefaultAvatar}
                />
              </Box>
            </HStack>
          </MenuButton>
          <MenuList
            bg={useColorModeValue('white', 'gray.900')}
            borderColor={useColorModeValue('gray.200', 'gray.700')}>
            <MenuItem onClick={() => navigate('/dashboard/profile')}>
              <Trans id="Profile">Profile</Trans>
            </MenuItem>
            {canCreateOrganization() ? (
              <MenuItem as={Link} to={APP_PATH.SWITCH} data-testid="switch_accounts">
                <Trans id="Switch Accounts">Switch accounts</Trans>
              </MenuItem>
            ) : null}
            <MenuDivider />
            <MenuItem onClick={handleLogout}>
              <Trans id="Sign out">Sign out</Trans>
            </MenuItem>
          </MenuList>
        </Menu>
      </HStack>
    </Flex>
  );
};

export default MobileNav;
