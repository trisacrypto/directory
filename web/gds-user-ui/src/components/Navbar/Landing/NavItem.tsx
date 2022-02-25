import React from "react";
import { LinkProps, NavLink } from "react-router-dom";
import { LinkBox, BoxProps, Text } from "@chakra-ui/react";

type Props = LinkProps &
  BoxProps & {
    pageName: string;
    to: string;
    disabled?: boolean;
  };

export const NavItem: React.FC<Props> = ({
  to,
  pageName,
  disabled,
  ...props
}): React.ReactElement => {
  if (disabled) {
    return (
      <Text padding={5} fontSize={16} color="gray.500" fontWeight="700">
        {pageName}
      </Text>
    );
  }
  return (
    <LinkBox
      as={NavLink}
      to={to}
      padding={5}
      fontSize={16}
      width="50%"
      color="gray.400"
      fontWeight="700"
      activeStyle={{
        color: "black",
        borderBottom: "5px solid gold",
      }}
      isActive={(match: unknown) => !!match}
      {...props}
    >
      {pageName}
    </LinkBox>
  );
};
