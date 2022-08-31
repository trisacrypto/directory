import {
  Icon,
  Box,
  ComponentWithAs,
  IconProps,
  styled,
  ListItem,
  ListItemProps
} from '@chakra-ui/react';
import { ReactNode } from 'react';
import { IconType } from 'react-icons';
import { NavLink as RouterLink, useLocation } from 'react-router-dom';
import ArrowIcon from './ArrowIcon';

export const StyledNavItem = styled(ListItem, {
  baseStyle: ({ isActive, isSubMenu }: any) => {
    return {
      w: '100%',
      py: '15px',
      cursor: 'pointer',
      position: 'relative',
      textDecor: 'none',
      color: '#A4A6B3',
      alignItems: 'center',
      display: 'flex',
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
      },
      '& > svg': {
        verticalAlign: 'text-bottom'
      },
      ...(isActive && {
        borderLeft: '2px solid #DDE2FF',
        background: 'hsla(231, 12%, 66%, 0.16)',
        width: '100%',
        color: '#DDE2FF'
      }),
      ...(isSubMenu && {
        pl: 10
      })
    };
  }
});

export interface NavItemProps extends ListItemProps {
  icon?: IconType | ComponentWithAs<'svg', IconProps>;
  href?: string;
  children: ReactNode;
  path?: string;
  hasChildren?: boolean;
  onOpen?: () => void;
  isCollapse?: boolean;
  isSubMenu?: boolean;
}

const NavItem = ({
  icon,
  children,
  hasChildren,
  href = '#',
  onOpen,
  path,
  isCollapse,
  ...rest
}: NavItemProps) => {
  const location = useLocation();
  const isActive = location.pathname === path;

  if (hasChildren) {
    return (
      <StyledNavItem
        w="100%"
        to={'/#'}
        onClick={onOpen}
        display="flex"
        alignItems="center"
        justifyContent="space-between"
        pr={3}
        isActive={isActive}
        {...rest}>
        <Box>
          {icon && (
            <Icon
              mr="4"
              _groupHover={{
                color: 'white'
              }}
              as={icon}
            />
          )}
          {children}
        </Box>
        <ArrowIcon open={isCollapse!} />
      </StyledNavItem>
    );
  }

  return (
    <StyledNavItem w="100%" as={RouterLink} display="flex" to={href} isActive={isActive} {...rest}>
      {icon && (
        <Icon
          mr="4"
          _groupHover={{
            color: 'white'
          }}
          as={icon}
        />
      )}
      {children}
    </StyledNavItem>
  );
};

export default NavItem;
