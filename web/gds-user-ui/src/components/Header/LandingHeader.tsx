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
  DrawerCloseButton
} from '@chakra-ui/react';
import { MenuIcon, CloseIcon } from '../Icon';
import Logo from 'components/ui/Logo';
import MenuItem from 'components/Menu/Landing/MenuItem';
import { colors } from 'utils/theme';
import LanguagesDropdown from 'components/LanguagesDropdown';

const LandingHeader = (props: FlexProps): JSX.Element => {
  const [show, setShow] = React.useState(false);
  const iconColor = useColorModeValue('black', 'white');
  const { isOpen, onOpen, onClose } = useDisclosure();

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
              <Link href="/" _active={{ outline: 'none' }} _focus={{ outline: 'none' }}>
                <Logo w={{ base: '50px', md: '100px' }} color={['colors.system.blue']} />
              </Link>
            </Box>
            <Box ml="auto" display={{ base: 'block', sm: 'none' }} onClick={onOpen}>
              {show ? <CloseIcon color={iconColor} /> : <MenuIcon color={iconColor} />}
            </Box>

            <Stack ml="auto" display={{ base: 'none', sm: 'flex' }} direction={['column', 'row']}>
              <MenuItem to="/#about">About TRISA</MenuItem>
              <MenuItem to="https://trisa.dev">Documentation </MenuItem>
              {/* <MenuItem isLast to="/auth/login">
                Login{' '}
              </MenuItem> */}
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
                    <MenuItem to="/#about">About TRISA </MenuItem>
                    <MenuItem to="https://trisa.dev">Documentation </MenuItem>
                    <MenuItem to="/auth/login">Login </MenuItem>
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
