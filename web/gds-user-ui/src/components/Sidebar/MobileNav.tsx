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
  useDisclosure
} from '@chakra-ui/react';
import { FiMenu } from 'react-icons/fi';
import LanguagesDropdown from 'components/LanguagesDropdown';
import { useDispatch, useSelector } from 'react-redux';
import { clearCookies } from 'utils/cookies';
import { useNavigate } from 'react-router-dom';
import DefaultAvatar from 'assets/default_avatar.svg';
import { resetStore } from 'application/store';
import { userSelector, logout } from 'modules/auth/login/user.slice';
import { Trans } from '@lingui/react';
import ChooseAnAccount from 'components/ChooseAnAccount';

interface MobileProps extends FlexProps {
  onOpen: () => void;
  isLoading?: boolean;
}
const MobileNav = ({ onOpen, ...rest }: MobileProps) => {
  const {
    isOpen: isAccountSwitchOpen,
    onOpen: onAccountSwitchOpen,
    onClose: onAccountSwitchClose
  } = useDisclosure();

  const dispatch = useDispatch();
  const { user } = useSelector(userSelector);
  const navigate = useNavigate();
  const handleLogout = (e: any) => {
    e.preventDefault();
    clearCookies();
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
            <MenuItem onClick={onAccountSwitchOpen}>
              <Trans id="Switch Accounts">Switch accounts</Trans>
              <ChooseAnAccount isOpen={isAccountSwitchOpen} onClose={onAccountSwitchClose} />
            </MenuItem>
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
