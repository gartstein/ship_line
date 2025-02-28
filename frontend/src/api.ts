// src/api.ts
const API_BASE_URL = "http://localhost:8080/v1";

// ✅ Fetch pack sizes as an array
export const fetchPackSizes = async () => {
    const response = await fetch(`${API_BASE_URL}/pack-sizes`);
    if (!response.ok) throw new Error(`Error ${response.status}: Failed to fetch pack sizes.`);

    const data = await response.json();
    return { pack_sizes: Array.isArray(data.pack_sizes) ? data.pack_sizes : [] }; // Ensure it's an array
};

// ✅ Send pack sizes as a JSON array with key `pack_sizes`
export const updatePackSizes = async (packSizes: number[]) => {
    const response = await fetch(`${API_BASE_URL}/pack-sizes`, {
        method: "PUT",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ pack_sizes: packSizes }), // ✅ Use `pack_sizes`
    });

    if (!response.ok) throw new Error(`Error ${response.status}: Failed to update pack sizes.`);
    return response.json();
};

// ✅ Use correct API endpoint for pack calculation
export const calculatePacks = async (items: number) => {
    const response = await fetch(`${API_BASE_URL}/calc?items=${items}`);

    if (!response.ok) {
        throw new Error(`Error ${response.status}: Failed to calculate pack sizes.`);
    }

    const data = await response.json();

    // ✅ Convert `packsUsed` object into an array [{ pack: 5000, count: 6 }]
    const formattedResult = Object.entries(data.packsUsed).map(([pack, count]) => ({
        pack: parseInt(pack, 10),
        count: count as number,
    }));

    return {
        itemsOrdered: data.itemsOrdered,
        totalItemsUsed: data.totalItemsUsed,
        result: formattedResult, // ✅ Ensure result is an array
    };
};

export const deletePackSize = async (size: number) => {
    const response = await fetch(`${API_BASE_URL}/pack-sizes/${size}`, {
        method: "DELETE",
    });
    if (!response.ok) {
        throw new Error("Failed to delete pack size");
    }
    return response.json();
};