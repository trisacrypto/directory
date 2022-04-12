import React from 'react';
import { Box, Flex, FlexProps, useColorModeValue, Link, Container, HStack } from '@chakra-ui/react';
import { MenuIcon, CloseIcon } from '../Icon';
import Logo from 'components/ui/Logo';
import MenuItem from 'components/Menu/Landing/MenuItem';
import { colors } from 'utils/theme';

const LandingHeader = (props: FlexProps): JSX.Element => {
  const [show, setShow] = React.useState(false);
  const toggleMenu = () => setShow(!show);
  const iconColor = useColorModeValue('black', 'white');
  return (
    <Flex
      width="100%"
      position={'relative'}
      p={8}
      bg={'transparent'}
      color={colors.system.blue}
      {...props}>
      <Container maxW={'5xl'}>
        <Box display={{ base: 'block', md: 'none' }} onClick={toggleMenu}>
          {show ? <CloseIcon color={iconColor} /> : <MenuIcon color={iconColor} />}
        </Box>

        <Box
          display={{ base: show ? 'block' : 'none', md: 'block' }}
          flexBasis={{ base: '100%', md: 'auto' }}>
          <Flex
            align="center"
            justify={['center', 'space-between', 'space-between', 'space-between']}
            direction={['column', 'row', 'row', 'row']}
            pt={[4, 4, 0, 0]}>
            <Box>
              <Link href="/">
                <Logo w="100px" color={['colors.system.blue']} />
              </Link>
            </Box>
            <HStack>
              <MenuItem to="/#about">About TRISA </MenuItem>
              {/* <MenuItem to="/about">About Us</MenuItem> */}
              <MenuItem to="https://trisa.dev">Documentation </MenuItem>
              {/* <MenuItem to="/login" isLast>
            Login
          </MenuItem> */}
            </HStack>
          </Flex>
        </Box>
      </Container>
    </Flex>
  );
};

export default LandingHeader;
