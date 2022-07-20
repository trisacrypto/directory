import React from 'react';
import {
  DrawerBody,
  Box,
  Flex,
  FlexProps,
  useColorModeValue,
  Link,
  Container,
  Stack,
  VStack,
  Drawer,
  DrawerOverlay,
  DrawerContent,
  DrawerHeader,
  useDisclosure,
  DrawerCloseButton,
  Text,
  Button
} from '@chakra-ui/react';
import { MenuIcon, CloseIcon } from '../Icon';
import Logo from 'components/ui/Logo';
import MenuItem from 'components/Menu/Landing/MenuItem';
import { colors } from 'utils/theme';
import { Trans } from '@lingui/react';
import LanguagesDropdown from 'components/LanguagesDropdown';
import { NavLink, Outlet } from 'react-router-dom';
import { useLanguageProvider } from 'contexts/LanguageContext';
import { TRISA_BASE_URL } from 'constants/trisa-base-url';

const LandingHeader = (props: FlexProps): JSX.Element => {
  const [show, setShow] = React.useState(false);
  const iconColor = useColorModeValue('black', 'white');
  const { isOpen, onOpen, onClose } = useDisclosure();
  const [locale] = useLanguageProvider();

  return (
    <Flex
      width="100%"
      position={'relative'}
      p={{ base: 4, md: 8 }}
      bg={'transparent'}
      color={colors.system.blue}
      {...props}>
      <Container maxW={'5xl'}>
        <Box flexBasis={{ base: '100%', md: 'auto' }}>
          <Flex align="center" justify={{ md: 'space-between' }}>
            <Box>
              <NavLink to={'/'}>
                <Link _active={{ outline: 'none' }} _focus={{ outline: 'none' }}>
                  <Logo w={{ base: '50px', md: '100px' }} color={['colors.system.blue']} />
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
              <MenuItem to="/#about">
                <Trans id="About TRISA">About TRISA</Trans>
              </MenuItem>
              <MenuItem data-testid="documentation" to={`${TRISA_BASE_URL}/${locale}`}>
                <Trans id="Documentation">Documentation</Trans>
              </MenuItem>
              <Stack>
                <NavLink to={'/auth/login'}>
                  <Button variant="secondary">
                    <Trans id="Login">Login</Trans>
                  </Button>
                </NavLink>
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
                <DrawerBody mt="50px" px={0}>
                  <VStack
                    alignItems="start"
                    sx={{
                      p: {
                        color: '#fff',
                        paddingY: 2,
                        m: '0 !important',
                        w: '100%',
                        pl: '25px'
                      }
                    }}>
                    <MenuItem to="/#about">
                      <Trans id="About TRISA">About TRISA</Trans>{' '}
                    </MenuItem>
                    <MenuItem to={`${TRISA_BASE_URL}/${locale}`}>
                      <Trans id="Documentation">Documentation</Trans>
                    </MenuItem>
                    <MenuItem to="/auth/login">
                      <Trans id="Login">Login</Trans>
                    </MenuItem>
                  </VStack>
                </DrawerBody>
              </DrawerContent>
            </Drawer>
          </Flex>
        </Box>
      </Container>
    </Flex>
  );
};

export default LandingHeader;
