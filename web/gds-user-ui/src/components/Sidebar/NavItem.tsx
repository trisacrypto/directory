import { Flex, FlexProps, Icon, Link, Box, Text } from '@chakra-ui/react';
import { ReactText } from 'react';
import { IconType } from 'react-icons';
import { Link as RouterLink } from 'react-router-dom';
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
      width: '260px',
      height: '100%',
      top: 0,
      color: 'white',
      left: 0,
      right: 0,
      borderLeft: 2,
      borderLeftStyle: 'solid',
      borderLeftColor: '#DDE2FF'
    }
  }
});

const NavItem = ({ icon, children, href = '#', selected, ...rest }: NavItemProps) => {
  return (
    <RouterLink to={href}>
      <Flex
        align="center"
        borderRadius="md"
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
    </RouterLink>
  );
};

export default NavItem;
