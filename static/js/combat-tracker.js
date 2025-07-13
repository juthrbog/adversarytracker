// Combat Tracker for Daggerheart Adversary Tracker
document.addEventListener('alpine:init', () => {
    Alpine.data('combatTracker', () => ({
        combatStarted: false,
        currentTurnIndex: -1,
        combatants: [],
        
        startCombat() {
            if (this.combatStarted) return;
            
            this.combatStarted = true;
            this.setupCombatants();
            this.sortByInitiative();
            this.renderCombatants();
            this.nextTurn();
            
            document.getElementById('combat-tracker').classList.remove('hidden');
            document.getElementById('start-combat').disabled = true;
            document.getElementById('next-turn').disabled = false;
            document.getElementById('end-combat').disabled = false;
        },
        
        setupCombatants() {
            // This will be populated by the Go template
            // The template will inject JavaScript to add combatants
        },
        
        sortByInitiative() {
            this.combatants.sort((a, b) => b.initiative - a.initiative);
        },
        
        renderCombatants() {
            const tbody = document.getElementById('initiative-order');
            tbody.innerHTML = '';
            
            this.combatants.forEach((combatant, index) => {
                const row = document.createElement('tr');
                row.id = combatant.id;
                row.className = index === this.currentTurnIndex ? 'bg-yellow-100' : '';
                
                row.innerHTML = `
                    <td class="py-2 px-4 border-b border-gray-200">${combatant.initiative}</td>
                    <td class="py-2 px-4 border-b border-gray-200 font-medium">${combatant.name}</td>
                    <td class="py-2 px-4 border-b border-gray-200">
                        <div class="flex items-center">
                            <button class="text-red-600 hover:text-red-800 mr-1" onclick="window.combatTracker.modifyHp('${combatant.id}', -1)">-</button>
                            <span>${combatant.currentHp}/${combatant.maxHp}</span>
                            <button class="text-green-600 hover:text-green-800 ml-1" onclick="window.combatTracker.modifyHp('${combatant.id}', 1)">+</button>
                        </div>
                    </td>
                    <td class="py-2 px-4 border-b border-gray-200">
                        <select class="text-sm border rounded" onchange="window.combatTracker.updateStatus('${combatant.id}', this.value)">
                            <option value="Normal" ${combatant.status === 'Normal' ? 'selected' : ''}>Normal</option>
                            <option value="Poisoned" ${combatant.status === 'Poisoned' ? 'selected' : ''}>Poisoned</option>
                            <option value="Stunned" ${combatant.status === 'Stunned' ? 'selected' : ''}>Stunned</option>
                            <option value="Unconscious" ${combatant.status === 'Unconscious' ? 'selected' : ''}>Unconscious</option>
                        </select>
                    </td>
                    <td class="py-2 px-4 border-b border-gray-200">
                        <button class="text-red-600 hover:text-red-800" onclick="window.combatTracker.removeCombatant('${combatant.id}')">Remove</button>
                    </td>
                `;
                
                tbody.appendChild(row);
            });
        },
        
        nextTurn() {
            if (!this.combatStarted || this.combatants.length === 0) return;
            
            this.currentTurnIndex = (this.currentTurnIndex + 1) % this.combatants.length;
            document.getElementById('current-turn').textContent = this.combatants[this.currentTurnIndex].name;
            this.renderCombatants();
        },
        
        endCombat() {
            this.combatStarted = false;
            this.currentTurnIndex = -1;
            this.combatants = [];
            
            document.getElementById('combat-tracker').classList.add('hidden');
            document.getElementById('start-combat').disabled = false;
            document.getElementById('next-turn').disabled = true;
            document.getElementById('end-combat').disabled = true;
            document.getElementById('current-turn').textContent = '-';
            document.getElementById('initiative-order').innerHTML = '';
        },
        
        modifyHp(id, amount) {
            const combatant = this.combatants.find(c => c.id === id);
            if (combatant) {
                combatant.currentHp = Math.max(0, Math.min(combatant.maxHp, combatant.currentHp + amount));
                this.renderCombatants();
            }
        },
        
        updateStatus(id, status) {
            const combatant = this.combatants.find(c => c.id === id);
            if (combatant) {
                combatant.status = status;
            }
        },
        
        removeCombatant(id) {
            const index = this.combatants.findIndex(c => c.id === id);
            if (index !== -1) {
                this.combatants.splice(index, 1);
                
                if (this.combatants.length === 0) {
                    this.endCombat();
                    return;
                }
                
                if (index <= this.currentTurnIndex) {
                    this.currentTurnIndex = Math.max(0, this.currentTurnIndex - 1);
                }
                
                this.renderCombatants();
                document.getElementById('current-turn').textContent = this.combatants[this.currentTurnIndex].name;
            }
        }
    }));
});

// Initialize and expose the combat tracker
document.addEventListener('DOMContentLoaded', () => {
    window.combatTracker = Alpine.data('combatTracker')();
    
    document.getElementById('start-combat').addEventListener('click', () => {
        window.combatTracker.startCombat();
    });
    
    document.getElementById('next-turn').addEventListener('click', () => {
        window.combatTracker.nextTurn();
    });
    
    document.getElementById('end-combat').addEventListener('click', () => {
        if (confirm('Are you sure you want to end combat?')) {
            window.combatTracker.endCombat();
        }
    });
});
