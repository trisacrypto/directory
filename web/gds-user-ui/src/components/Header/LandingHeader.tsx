import React from 'react';
import {
  DrawerBody,
  Box,
  Flex,
  FlexProps,
  useColorModeValue,
  Container,
  Stack,
  Drawer,
  DrawerOverlay,
  DrawerContent,
  useDisclosure,
  DrawerCloseButton,
  Button,
  Link
} from '@chakra-ui/react';
import { MenuIcon, CloseIcon } from '../Icon';
import MenuItem from 'components/Menu/Landing/MenuItem';
import { colors } from 'utils/theme';
import { Trans } from '@lingui/react';
import LanguagesDropdown from 'components/LanguagesDropdown';
import { NavLink } from 'react-router-dom';
import { useLanguageProvider } from 'contexts/LanguageContext';
import { TRISA_BASE_URL } from 'constants/trisa-base-url';
import useAuth from 'hooks/useAuth';
import CkLazyLoadImage from 'components/LazyImage';
import TrisaLogo from 'assets/TRISA-GDS-black.png';

const LandingHeader = (props: FlexProps): JSX.Element => {
  const [show] = React.useState(false);
  const iconColor = useColorModeValue('black', 'white');
  const { isOpen, onOpen, onClose } = useDisclosure();
  const [locale] = useLanguageProvider();
  const { isAuthenticated } = useAuth();
  const isLoggedIn = isAuthenticated();

  return (
    <Flex
      width="100%"
      position={'relative'}
      bg={'transparent'}
      boxShadow="md"
      color={colors.system.blue}
      {...props}>
      <Container maxW={'5xl'}>
        <Box flexBasis={{ base: '100%', md: 'auto' }}>
          <Flex align="center" justify={{ md: 'space-between' }}>
            <Box>
              <NavLink to={'/'}>
                <Link _active={{ outline: 'none' }} _focus={{ outline: 'none' }}>
                  <CkLazyLoadImage
                    src={TrisaLogo}
                    alt="Trisa logo"
                    objectFit="cover"
                    height="100px"
                    transform="translateX(-38px)"
                    sx={{ aspectRatio: '2/1' }}
                  />
                </Link>
              </NavLink>
            </Box>
            <Box ml="auto" display={{ base: 'block', sm: 'none' }} onClick={onOpen}>
              {show ? <CloseIcon color={iconColor} /> : <MenuIcon color={iconColor} />}
            </Box>

            <Stack
              isInline
              align="center"
              justify="flex-end"
              ml={{ base: 'auto', md: 0 }}
              alignItems={'center'}
              display={{ base: 'none', sm: 'flex' }}
              direction={['column', 'row']}>
              <Stack pr={2}>
                <LanguagesDropdown />
              </Stack>
              <MenuItem to="https://travelrule.io" data-testid="about">
                <Trans id="About TRISA">About TRISA</Trans>
              </MenuItem>
              <MenuItem data-testid="documentation" to={`${TRISA_BASE_URL}/${locale}`}>
                <Trans id="Documentation">Documentation</Trans>
              </MenuItem>
              <Stack>
                {!isLoggedIn ? (
                  <NavLink to={'/auth/login'}>
                    <Button variant="secondary" data-cy="nav-login-bttn">
                      <Trans id="Login">Login</Trans>
                    </Button>
                  </NavLink>
                ) : (
                  <NavLink to={'/dashboard/overview'}>
                    <Button variant="secondary">
                      <Trans id="Your dashboard">Your dashboard</Trans>
                    </Button>
                  </NavLink>
                )}
              </Stack>
            </Stack>

            {/* mobile drawer */}
            <Drawer placement="right" onClose={onClose} isOpen={isOpen} size="xs">
              <DrawerOverlay />
              <DrawerContent bg="#262626">
                <DrawerCloseButton
                  left={'15px'}
                  color="#fff"
                  sx={{
                    '.chakra-icon path': {
                      fill: '#fff'
                    }
                  }}
                />
                <DrawerBody mt="50px" px={5}>
                  <MenuItem to="https://travelrule.io" color="white" pb={0}>
                    <Trans id="About TRISA">About TRISA</Trans>
                  </MenuItem>
                  <MenuItem to={`${TRISA_BASE_URL}/${locale}`} color="white">
                    <Trans id="Documentation">Documentation</Trans>
                  </MenuItem>
                  <Stack>
                    {!isLoggedIn ? (
                      <NavLink to={'/auth/login'}>
                        <Button variant="secondary">
                          <Trans id="Login">Login</Trans>
                        </Button>
                      </NavLink>
                    ) : (
                      <NavLink to={'/dashboard/overview'}>
                        <Button variant="secondary">
                          <Trans id="Your dashboard">Your dashboard</Trans>
                        </Button>
                      </NavLink>
                    )}
                  </Stack>
                </DrawerBody>
              </DrawerContent>
            </Drawer>
          </Flex>
        </Box>
      </Container>
    </Flex>
  );
};

export default React.memo(LandingHeader);
