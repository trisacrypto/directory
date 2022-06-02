import { Flex, FlexProps, Icon, Link } from '@chakra-ui/react';
import { ReactText } from 'react';
import { IconType } from 'react-icons';

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
      width: '100%',
      height: '100%',
      top: 0,
      left: 0,
      borderLeft: 2,
      borderLeftStyle: 'solid',
      borderLeftColor: '#DDE2FF'
    }
  }
});

const NavItem = ({ icon, children, href = '#', selected, ...rest }: NavItemProps) => {
  return (
    <Link href={href} {...getLinkStyle()}>
      <Flex
        align="center"
        borderRadius="md"
        role="group"
        color={selected ? 'white' : '#8391a2'}
        fontSize="0.9375rem"
        _hover={{
          color: 'white'
        }}
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
        {children}
      </Flex>
    </Link>
  );
};

export default NavItem;
