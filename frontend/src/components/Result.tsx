import React from "react";
import { Typography, Paper } from "@mui/material";

interface ResultProps {
    result: {
        itemsOrdered: number;
        totalItemsUsed: number;
        packsUsed: Record<string, number>;
    } | null;
}

const Result: React.FC<ResultProps> = ({ result }) => {
    if (!result) {
        return <Typography color="textSecondary">Enter items and click "Calculate" to see results.</Typography>;
    }

    return (
        <Paper elevation={3} style={{ padding: "1rem", marginTop: "1rem" }}>
            <Typography variant="h6">Results</Typography>
            <Typography><strong>Items Ordered:</strong> {result.itemsOrdered}</Typography>
            <Typography><strong>Total Items Used:</strong> {result.totalItemsUsed}</Typography>
            <Typography><strong>Packs Used:</strong></Typography>
            <ul>
                {Object.entries(result.packsUsed).map(([packSize, count]) => (
                    <li key={packSize}>{`${count} × ${packSize}`}</li>
                ))}
            </ul>
        </Paper>
    );
};

export default Result;