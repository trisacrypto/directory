import React from 'react';
import { Link, Box, Flex, Text, Button, Stack } from "@chakra-ui/react";
import {MenuIcon , CloseIcon} from '../Icon/index'
import Logo from "../UI/Logo";
import MenuItem from "../Menu/Landing/MenuItem";

interface IProps {

}

const LandingHeader = (props : IProps) => {
  const [show, setShow] = React.useState(false);
  const toggleMenu = () => setShow(!show);

  return (
    <Flex
      as="nav"
      align="center"
      justify="space-between"
      wrap="wrap"
      w="100%"
      mb={8}
      p={8}
      bg={["system.link", "system.link", "transparent", "transparent"]}
      color={["white", "white", "system.link", "system.link"]}
      {...props}
    >
      <Flex align="center">
        <Logo
          w="100px"
          color={["white", "white", "system.link", "system.link"]}
        />
      </Flex>

      <Box display={{ base: "block", md: "none" }} onClick={toggleMenu}>
        {show ? <CloseIcon /> : <MenuIcon />}
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
          <MenuItem to="/doc">Documentation </MenuItem>
          <MenuItem to="/login" isLast>
            Login
          </MenuItem>
        </Flex>
      </Box>
    </Flex>
  );
};

export default LandingHeader;