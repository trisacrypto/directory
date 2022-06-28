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
  Heading
} from '@chakra-ui/react';
import trisaLogo from '../../assets/trisa.svg';
import NavItem from './NavItem';
import MenuItems from '../../utils/menu';
import { MdContactSupport } from 'react-icons/md';
import { IoLogoSlack } from 'react-icons/io';

interface SidebarProps extends BoxProps {
  onClose: () => void;
}

const SidebarContent = ({ onClose, ...rest }: SidebarProps) => {
  return (
    <Box
      transition="3s ease"
      bg={useColorModeValue('white', 'gray.900')}
      borderRight="1px"
      borderRightColor={useColorModeValue('gray.200', 'gray.700')}
      w={{ base: 'full', md: 275 }}
      pos="fixed"
      px={2}
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
        {MenuItems.filter((m) => m.activated).map((menu) => (
          <NavItem key={menu.title} icon={menu.icon} href={menu.path || '/#'} path={menu.path}>
            {menu.title}
          </NavItem>
        ))}

        <Divider maxW="80%" my="16px !important" mx="auto !important" />
        <NavItem href="mailto:support@trisa.io" icon={MdContactSupport} w={'100%'}>
          Support
        </NavItem>
        <NavItem icon={IoLogoSlack} href="/#" w={'100%'}>
          Slack
        </NavItem>
      </VStack>
    </Box>
  );
};

export default SidebarContent;
