import {
  Box,
  BoxProps,
  Flex,
  useColorModeValue,
  CloseButton,
  Divider,
  VStack,
  Stack,
  Link,
  Icon,
  Text,
  Collapse,
  List,
  HStack
} from '@chakra-ui/react';
import trisaLogo from 'assets/TRISA-GDS-white.png';
import NavItem, { StyledNavItem } from './NavItem';
import MenuItems from '../../utils/menu';
import { MdContactSupport } from 'react-icons/md';
import { IoLogoSlack } from 'react-icons/io';
import { Fragment, useState } from 'react';
import { Trans } from '@lingui/react';
import { LazyLoadImage } from 'react-lazy-load-image-component';
import ChakraRouterLink from 'components/ChakraRouterLink';
import HDivider from 'components/ui/HDivider';
// import Version from 'components/Footer/Version';
import useFetchAppVersion from 'hooks/useFetchAppVersion';
interface SidebarProps extends BoxProps {
  onClose: () => void;
}

const SidebarContent = ({ onClose, ...rest }: SidebarProps) => {
  const [open, setOpen] = useState(false);
  const { appGitVersion, bffAndGdsVersion } = useFetchAppVersion();
  return (
    <Box
      transition="3s ease"
      bg={useColorModeValue('white', 'gray.900')}
      borderRight="1px"
      borderRightColor={useColorModeValue('gray.200', 'gray.700')}
      w={{ base: 'full', md: 275 }}
      pos="fixed"
      h="full"
      {...rest}>
      <Flex h="20" alignItems="center" mx="8" my={2} justifyContent="space-between">
        <ChakraRouterLink to="/dashboard/overview">
          <Stack width="100%" direction={['row']} height="200px">
            <LazyLoadImage src={trisaLogo} alt="GDS UI" />
          </Stack>
        </ChakraRouterLink>
        <CloseButton display={{ base: 'flex', md: 'none' }} onClick={onClose} />
      </Flex>
      <VStack alignItems="flex-start" justifyContent="center" spacing={0}>
        <List w="100%">
          {MenuItems.filter((m) => m.activated).map((menu, index) => (
            <Fragment key={index}>
              <NavItem
                key={menu.title}
                icon={menu.icon}
                href={menu.path || '/#'}
                path={menu.path}
                hasChildren={!!menu.children?.length}
                onOpen={() => setOpen(!open)}
                isCollapse={open}>
                {menu.title}
              </NavItem>
              {menu.children?.length && (
                <Collapse in={open} style={{ width: '100%' }}>
                  {menu.children &&
                    menu.children
                      .filter((m) => m.activated)
                      .map((child) => (
                        <NavItem
                          key={child.title}
                          icon={child.icon}
                          href={child.path || '/#'}
                          path={child.path}
                          isCollapse={false}
                          isSubMenu={true}>
                          {child.title}
                        </NavItem>
                      ))}
                </Collapse>
              )}
            </Fragment>
          ))}
        </List>
        <Divider maxW="80%" my="16px !important" mx="auto !important" />
        <List w="100%">
          <StyledNavItem
            w="100%"
            display="flex"
            alignItems="center"
            color="#8391a2"
            role="group"
            href="mailto:support@rotational.io"
            as={Link}>
            <Icon
              mr="4"
              fontSize="16"
              _groupHover={{
                color: 'white'
              }}
              as={MdContactSupport}
            />
            <Text
              _groupHover={{
                color: 'white'
              }}>
              <Trans id="Support">Support</Trans>
            </Text>
          </StyledNavItem>
          <StyledNavItem
            href="https://trisa-workspace.slack.com/"
            w={'100%'}
            display="flex"
            alignItems="center"
            isExternal
            color="#8391a2"
            role="group"
            as={Link}>
            <Icon
              mr="4"
              fontSize="16"
              _groupHover={{
                color: 'white'
              }}
              as={IoLogoSlack}
            />
            <Text
              _groupHover={{
                color: 'white'
              }}>
              Slack
            </Text>
          </StyledNavItem>
        </List>
      </VStack>
      <HStack
        fontSize="0.6em"
        px={2}
        justifyContent="center"
        mx={'auto'}
        position={'absolute'}
        bottom={5}
        textAlign={'center'}
        color="white"
        width="100%">
        <Text>App: {appGitVersion || 'N/A'}</Text>
        <HDivider />
        <Text>BFF & GDS: {bffAndGdsVersion || 'N/A'}</Text>
      </HStack>
    </Box>
  );
};

export default SidebarContent;
