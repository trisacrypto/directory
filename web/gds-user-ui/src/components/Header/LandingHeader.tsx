import React from 'react';
import { Link, Box, Flex, FlexProps } from "@chakra-ui/react";
import {MenuIcon , CloseIcon} from '../icon'
import Logo from "../ui/Logo";
import MenuItem from "../menu/landing/MenuItem";
import { colors } from '../../utils/theme';

const LandingHeader = (props : FlexProps) : JSX.Element => {
  const [show, setShow] = React.useState(false);
  const toggleMenu = () => setShow(!show);

  return (
    <Flex
      as="nav"
      align="center"
      justify="space-between"
      wrap="wrap"
      w="100%"
      p={8}
      bg={["white", "white", "transparent", "transparent"]}
      color={colors.system.blue}
      {...props}
    >
      <Flex align="center">
        <Logo
          w="100px"
          color={["system.blue"]}
        />
      </Flex>

      <Box display={{ base: "block", md: "none" }} onClick={toggleMenu}>
        {show ? <CloseIcon  /> : <MenuIcon />}
      </Box>

      <Box
        display={{ base: show ? "block" : "none", md: "block" }}
        flexBasis={{ base: "100%", md: "auto" }}
      >
        <Flex
          align="center"
          justify={["center", "space-between", "flex-end", "flex-end"]}
          direction={["column", "row", "row", "row"]}
          pt={[4, 4, 0, 0]}
        >
          <MenuItem to="/">Home</MenuItem>
          <MenuItem to="/about">About Us</MenuItem>
          <MenuItem to="https://trisa.dev">Documentation </MenuItem>
          <MenuItem to="/login" isLast>
            Login
          </MenuItem>
        </Flex>
      </Box>
    </Flex>
  );
};

export default LandingHeader;