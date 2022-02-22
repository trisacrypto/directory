import { Link, Text  } from "@chakra-ui/react";

interface IProps {
    children: React.ReactNode;
    isLast?: boolean;
    to: string;
} 

const MenuItem = ({ children, isLast, to = "/", ...rest }: IProps) : JSX.Element => {
  return (
    <Text
      mb={{ base: isLast ? 0 : 8, sm: 0 }}
      mr={{ base: 0, sm: isLast ? 0 : 8 }}
      display="block"
      {...rest}
    >
      <Link href={to}>{children}</Link>
    </Text>
  );
};

export default MenuItem