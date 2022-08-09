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
  Show
} from '@chakra-ui/react';
import { FiMenu } from 'react-icons/fi';
import LanguagesDropdown from 'components/LanguagesDropdown';
import { useDispatch, useSelector } from 'react-redux';
import useCustomAuth0 from 'hooks/useCustomAuth0';
import { removeCookie } from 'utils/cookies';
import { useNavigate } from 'react-router-dom';
import DefaultAvatar from 'assets/default_avatar.svg';
import { resetStore } from 'application/store';
import Storage from 'reduxjs-toolkit-persist/lib/storage/session';
import AvatarContentLoader from 'components/ContentLoader/Avatar';
import { userSelector, logout } from 'modules/auth/login/user.slice';
import Logo from 'components/ui/Logo';

interface MobileProps extends FlexProps {
  onOpen: () => void;
  isLoading?: boolean;
}
const DEFAULT_AVARTAR = 'https://www.gravatar.com/avatar/205e460b479e2e5b48aec07710c08d50?s=200';
const MobileNav = ({ onOpen, ...rest }: MobileProps) => {
  const dispatch = useDispatch();
  const { user } = useSelector(userSelector);
  const navigate = useNavigate();
  const { auth0Logout } = useCustomAuth0();
  const handleLogout = (e: any) => {
    e.preventDefault();
    removeCookie('access_token');
    localStorage.removeItem('persist:root');
    dispatch(logout());
    resetStore();
    navigate('/');
  };
  return (
    <Flex
      ml={{ base: 0, md: 60 }}
      px={{ base: 4, md: 4 }}
      height="20"
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
      <Show below="md">
        <Logo sx={{ '& img': { width: '50%' } }} />
      </Show>
      <HStack
        spacing={{ base: '0', md: '6' }}
        w={{ base: '100%', md: 'none' }}
        gap={{ base: 4, md: 0 }}
        justifyContent="end">
        <HStack>
          <LanguagesDropdown />
        </HStack>
        <Divider orientation="vertical" height={8} />
        <Menu>
          <MenuButton transition="all 0.3s" _focus={{ boxShadow: 'none' }}>
            <HStack>
              <Show above="lg">
                <Text fontSize="sm" color="blackAlpha.700">
                  {user?.name || 'Guest'}
                </Text>
              </Show>
              <Box borderRadius="50%" borderWidth={2} padding={0.5}>
                <Avatar
                  size={'md'}
                  height="43.3"
                  w="43.3"
                  src={user?.pictureUrl || DefaultAvatar}
                />
              </Box>
            </HStack>
          </MenuButton>
          <MenuList
            bg={useColorModeValue('white', 'gray.900')}
            borderColor={useColorModeValue('gray.200', 'gray.700')}>
            <MenuDivider />
            <MenuItem onClick={handleLogout}>Sign out</MenuItem>
          </MenuList>
        </Menu>
      </HStack>
    </Flex>
  );
};

export default MobileNav;
