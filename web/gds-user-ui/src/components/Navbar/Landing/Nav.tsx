import React from 'react';
import { NavItem } from './NavItem';
import { HStack, Flex } from '@chakra-ui/react';
import Logo from '../../ui/Logo';
import { t } from '@lingui/macro';
export const NavBar = ({ ...props }): React.ReactElement => {
  return (
    <Flex
      direction="column"
      px={14}
      paddingTop={10}
      paddingBottom={6}
      bg="white"
      boxShadow="md"
      {...props}>
      <Flex px={4} justifyContent="space-between">
        <Logo />

        <HStack justify="center" alignItems="flex-start">
          <NavItem to="home" pageName={t`Home`} />
          <NavItem to="about" pageName={t`about us`} />
          <NavItem to="documentation" pageName={t`Documentation`} />
          <NavItem to="login" pageName={t`Log in`} />
        </HStack>
      </Flex>
    </Flex>
  );
};
