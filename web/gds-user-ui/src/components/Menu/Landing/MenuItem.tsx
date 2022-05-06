import { Link, Text } from '@chakra-ui/react';

interface MenuItemProps {
  children: React.ReactNode;
  isLast?: boolean;
  to: string;
}

const MenuItem = ({ children, isLast, to = '/', ...rest }: MenuItemProps): JSX.Element => {
  return (
    <Text
      mb={{ base: isLast ? 0 : 4, sm: 0 }}
      mr={{ base: 0, sm: isLast ? 8 : 0 }}
      pl={isLast ? 8 : 0}
      display="block"
      {...rest}>
      {to.startsWith('http') ? <a href={to}>{children}</a> : <Link href={to}>{children}</Link>}
    </Text>
  );
};

export default MenuItem;
