import React from 'react';
import { Box, Flex, FlexProps, useColorModeValue, Link } from '@chakra-ui/react';
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
      as="nav"
      align="center"
      justify="space-between"
      wrap="wrap"
      w="100%"
      p={8}
      bg={'transparent'}
      px={40}
      color={colors.system.blue}
      {...props}>
      <Flex align="center">
        <Link href="/">
          <Logo w="100px" color={['colors.system.blue']} />
        </Link>
      </Flex>

      <Box display={{ base: 'block', md: 'none' }} onClick={toggleMenu}>
        {show ? <CloseIcon color={iconColor} /> : <MenuIcon color={iconColor} />}
      </Box>

      <Box
        display={{ base: show ? 'block' : 'none', md: 'block' }}
        flexBasis={{ base: '100%', md: 'auto' }}>
        <Flex
          align="center"
          justify={['center', 'space-between', 'flex-end', 'flex-end']}
          direction={['column', 'row', 'row', 'row']}
          pt={[4, 4, 0, 0]}>
          <MenuItem to="/#about">About Trisa </MenuItem>
          {/* <MenuItem to="/about">About Us</MenuItem> */}
          <MenuItem to="https://trisa.dev">Documentation </MenuItem>
          {/* <MenuItem to="/login" isLast>
            Login
          </MenuItem> */}
        </Flex>
      </Box>
    </Flex>
  );
};

export default LandingHeader;
