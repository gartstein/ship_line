import React from "react";
import { render, fireEvent, screen } from "@testing-library/react";
import "@testing-library/jest-dom";
import Input from "../Input";

test("renders input field correctly", () => {
    render(<Input value="" onChange={() => {}} />);

    // Use getByLabelText instead of getByPlaceholderText
    const inputElement = screen.getByLabelText("Number of Items");
    expect(inputElement).toBeInTheDocument();
});

test("calls onChange handler when user types", () => {
    const handleChange = jest.fn();
    render(<Input value="" onChange={handleChange} />);

    const inputElement = screen.getByLabelText("Number of Items");

    fireEvent.change(inputElement, { target: { value: "100" } });
    expect(handleChange).toHaveBeenCalledTimes(1);
});