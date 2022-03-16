import { Meta, Story } from "@storybook/react";
import Sidebar from ".";

type SidebarProps = {
  children: React.ReactNode;
};

export default {
  title: "components/SideBar",
  component: Sidebar,
} as Meta<SidebarProps>;

const Template: Story<SidebarProps> = (args) => <Sidebar {...args} />;

export const Standard = Template.bind({});
Standard.args = {
  children: "This is my child component",
};
