import { Flex, FlexProps, Icon, Box, Text, chakra } from '@chakra-ui/react';
import { ReactText } from 'react';
import { IconType } from 'react-icons';
import { NavLink as RouterLink } from 'react-router-dom';

const ChakraRouterLink = chakra(RouterLink);
interface NavItemProps extends FlexProps {
  icon?: IconType;
  href?: string;
  children: ReactText;
  selected?: boolean;
}

const getLinkStyle: any = () => ({
  w: '100%',
  py: '0.9rem',
  cursor: 'pointer',
  position: 'relative',
  textDecor: 'none',
  pl: 7,
  _focus: { boxShadow: 'none' },
  _hover: {
    _after: {
      background: 'hsla(231, 12%, 66%, 0.16)',
      position: 'absolute',
      content: '""',
      height: '100%',
      top: 0,
      color: 'white',
      left: 0,
      right: 0,
      borderLeft: '2px solid #DDE2FF'
    }
  }
});

const getActiveLinkStyle = ({ isActive }: { isActive: boolean }) =>
  isActive
    ? {
        borderLeft: '2px solid #DDE2FF',
        background: 'hsla(231, 12%, 66%, 0.16)',
        width: '100%'
      }
    : {};

const NavItem = ({ icon, children, href = '#', selected, ...rest }: NavItemProps) => {
  return (
    <ChakraRouterLink w="100%" to={href} style={getActiveLinkStyle}>
      <Flex
        align="center"
        borderRadius="md"
        w="100%"
        role="group"
        color={selected ? 'white' : '#8391a2'}
        fontSize="0.9375rem"
        _hover={{
          color: 'white'
        }}
        {...getLinkStyle()}
        {...rest}>
        {icon && (
          <Icon
            mr="4"
            fontSize="16"
            _groupHover={{
              color: 'white'
            }}
            color={selected ? 'white' : '#8391a2'}
            as={icon}
          />
        )}
        <Box>
          <Text>{children}</Text>
        </Box>
      </Flex>
    </ChakraRouterLink>
  );
};

export default NavItem;
