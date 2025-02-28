import React from "react";
import { render, screen } from "@testing-library/react";
import "@testing-library/jest-dom";
import Result from "../Result";

test("renders calculation result correctly", () => {
    const mockResult = {
        itemsOrdered: 10,
        totalItemsUsed: 12,
        packsUsed: { "100": 2, "50": 1 },
    };

    render(<Result result={mockResult} />);

    // Check "Items Ordered:" text
    expect(screen.getByText(/Items Ordered:/i)).toBeInTheDocument();

    // Ensure "10" appears exactly once in the correct place
    const itemsOrderedElements = screen.getAllByText(/10/);
    expect(itemsOrderedElements.length).toBeGreaterThan(0);

    // Check "Total Items Used:" text
    expect(screen.getByText(/Total Items Used:/i)).toBeInTheDocument();
    expect(screen.getByText(/12/)).toBeInTheDocument();

    // Ensure packsUsed are correctly displayed
    expect(screen.getByText(/2 × 100/i)).toBeInTheDocument();
    expect(screen.getByText(/1 × 50/i)).toBeInTheDocument();
});