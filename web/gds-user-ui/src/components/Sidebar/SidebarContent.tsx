import {
  Box,
  BoxProps,
  Flex,
  useColorModeValue,
  Image,
  CloseButton,
  Divider,
  VStack,
  Stack,
  Heading,
  Link,
  Icon,
  Text,
  HStack,
  Collapse
} from '@chakra-ui/react';
import trisaLogo from '../../assets/trisa.svg';
import NavItem, { getLinkStyle, NavItemProps } from './NavItem';
import MenuItems from '../../utils/menu';
import { MdContactSupport } from 'react-icons/md';
import { IoLogoSlack } from 'react-icons/io';
import { FaChevronDown } from 'react-icons/fa';
import { useState } from 'react';

interface SidebarProps extends BoxProps {
  onClose: () => void;
}

const SubMenuItem = (props: NavItemProps) => <NavItem {...props} pl="3rem" />;

const SidebarContent = ({ onClose, ...rest }: SidebarProps) => {
  const [open, setOpen] = useState(false);
  return (
    <Box
      transition="3s ease"
      bg={useColorModeValue('white', 'gray.900')}
      borderRight="1px"
      borderRightColor={useColorModeValue('gray.200', 'gray.700')}
      w={{ base: 'full', md: 275 }}
      pos="fixed"
      // px={2}
      h="full"
      {...rest}>
      <Flex h="20" alignItems="center" mx="8" my={2} justifyContent="space-between">
        <Stack width="100%" direction={['row']}>
          <Image src={trisaLogo} alt="GDS UI" />
          <Heading size="sm" color="#FFFFFF" lineHeight={1.35}>
            Global Directory Service
          </Heading>
        </Stack>
        <CloseButton display={{ base: 'flex', md: 'none' }} onClick={onClose} />
      </Flex>
      <VStack alignItems="flex-start" justifyContent="center" spacing={0}>
        <VStack w="100%">
          {MenuItems.filter((m) => m.activated).map((menu) => (
            <>
              <NavItem
                key={menu.title}
                icon={menu.icon}
                href={menu.path || '/#'}
                path={menu.path}
                hasChildren={!!menu.children}
                onOpen={() => setOpen(!open)}>
                {menu.title}
                {menu.children && (
                  <FaChevronDown
                    style={{
                      transform: open ? 'rotate(180deg)' : undefined,
                      transition: '200ms'
                    }}
                  />
                )}
              </NavItem>
              {menu.children?.length && (
                <Collapse in={open} style={{ width: '100%' }}>
                  {menu.children &&
                    menu.children
                      .filter((m) => m.activated)
                      .map((child) => (
                        <SubMenuItem
                          key={child.title}
                          icon={child.icon}
                          href={child.path || '/#'}
                          path={child.path}>
                          {child.title}
                        </SubMenuItem>
                      ))}
                </Collapse>
              )}
            </>
          ))}
        </VStack>
        <Divider maxW="80%" my="16px !important" mx="auto !important" />
        <VStack w="100%">
          <Link
            w="100%"
            display="flex"
            alignItems="center"
            color="#8391a2"
            role="group"
            href="mailto:support@trisa.io"
            isExternal
            {...getLinkStyle()}>
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
              Support
            </Text>
          </Link>
          <Link
            href="https://trisa-workspace.slack.com/"
            w={'100%'}
            display="flex"
            alignItems="center"
            color="#8391a2"
            role="group"
            isExternal
            {...getLinkStyle()}>
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
          </Link>
          {/* <NavItem icon={IoLogoSlack} href="https://trisa-workspace.slack.com/" w={'100%'}>
          Slack
        </NavItem> */}
        </VStack>
      </VStack>
    </Box>
  );
};

export default SidebarContent;
